package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
	"time"
)

const jsonPath = "playlist.json"

func (p *Playlist) CreateSong(s *Song) {
	newSong := &Song{Name: s.Name, Duration: s.Duration, ID: rand.Int()}
	if s.ID != 0 {
		newSong.ID = s.ID
	}
	if p.Songs.Head == nil {
		p.Songs.Head = newSong
		p.Songs.tail = newSong
	} else {
		p.Songs.tail.next = newSong
		newSong.prev = p.Songs.tail
		p.Songs.tail = newSong
		newSong.next = p.Songs.Head
		p.Songs.Head.prev = newSong
	}
	if p.Songs.Current == nil {
		p.Songs.Current = newSong
	}

	fmt.Printf("Песня %s была добавлена в плейлист\n", newSong.Name)
}
func (p *Playlist) WriteSongToJson(s *Song) error {
	jsonMap := map[string]string{"Name": s.Name, "Duration": strconv.Itoa(int(s.Duration)), "ID": strconv.Itoa(rand.Int())}

	jsonData, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		log.Fatal(err)
	}

	// Decode JSON data into a slice of maps
	var data []map[string]string
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		log.Fatal(err)
	}

	// Append the new map to the slice of maps
	data = append(data, jsonMap)

	// Encode updated slice of maps back into JSON
	jsonData, err = json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	// Write updated JSON data to file
	err = ioutil.WriteFile(jsonPath, jsonData, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func (p *Playlist) LoadSongs() {

	err, data := UnmarshalJson()

	var songs []Song

	// iterate over the slice of maps and create Song structs
	for _, item := range data {

		song := Song{}
		if id, ok := item["ID"]; ok {
			song.ID, err = strconv.Atoi(id)
			if err != nil {
				panic(err)
			}
		}

		if name, ok := item["Name"]; ok {
			song.Name = name
		}
		if duration, ok := item["Duration"]; ok {
			durInt, err := strconv.Atoi(duration)
			if err != nil {
				panic(err)
			}
			song.Duration = time.Duration(durInt) * time.Second
		}
		songs = append(songs, song)
	}
	for _, song := range songs {
		p.CreateSong(&song)
	}

}

func (p *Playlist) ShowSong(s *Song) (map[string]string, error) {
	_, data := UnmarshalJson()

	for _, item := range data {

		itemID, err := strconv.Atoi(item["ID"])
		if err != nil {
			panic(err)
		}

		if itemID == s.ID {
			return item, nil
		}
	}

	return map[string]string{}, errors.New("Песня не найдена")
}

func (p *Playlist) UpdateSong(s *Song) *Song {
	_, data := UnmarshalJson()

	for _, item := range data {

		itemID, err := strconv.Atoi(item["ID"])
		if err != nil {
			panic(err)
		}

		if itemID == s.ID {
			updatedSong := Song{ID: itemID, Name: s.Name, Duration: s.Duration}

			for node := p.Songs.Head; node != nil; node = node.next {
				if node.ID == itemID {
					node.Name = s.Name
					node.Duration = s.Duration * time.Second

					p.UpdateJson(&updatedSong)
					return &updatedSong
				}
			}
		}
	}

	return s
}

func (p *Playlist) UpdateJson(s *Song) {
	err, data := UnmarshalJson()

	for _, item := range data {

		itemID, err := strconv.Atoi(item["ID"])
		if err != nil {
			panic(err)
		}

		if itemID == s.ID {
			item["Name"] = s.Name
			item["Duration"] = strconv.Itoa(int(s.Duration))
		}
	}

	newjsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	// Write updated JSON data to file
	err = ioutil.WriteFile(jsonPath, newjsonData, 0644)
	if err != nil {
		log.Fatal(err)
	}

}

func (p *Playlist) DeleteSong(s *Song) error {
	if p.Songs.Current.ID == s.ID && p.IsPlay {
		return errors.New("Данная песня сейчас проигрывается, пожалуйста попробуйте позже")
	}

	_, data := UnmarshalJson()

	for _, item := range data {

		itemID, err := strconv.Atoi(item["ID"])
		if err != nil {
			panic(err)
		}

		if itemID == s.ID {
			songForDelete := Song{ID: itemID, Name: s.Name, Duration: s.Duration}
			p.DelSongFromJson(&songForDelete)
			curr := p.Songs.Head
			for curr != nil && curr.ID != itemID {
				curr = curr.next
			}
			if curr == nil {
				return nil
			}

			if curr.prev == nil {
				p.Songs.Head = curr.next
				if p.Songs.Head != nil {
					p.Songs.Head.prev = nil
				}
			} else {

				curr.prev.next = curr.next
			}

			if curr.next == nil {
				p.Songs.tail = curr.prev
				if p.Songs.tail != nil {
					p.Songs.tail.next = nil
				}
			} else {

				curr.next.prev = curr.prev
			}
			return nil

		}
	}
	return errors.New("Песня не найдена")
}

func (p *Playlist) DelSongFromJson(s *Song) {
	_, data := UnmarshalJson()

	for index, item := range data {

		itemID, err := strconv.Atoi(item["ID"])
		if err != nil {
			panic(err)
		}

		if itemID == s.ID {
			data = append(data[:index], data[index+1:]...)
		}
	}

	newjsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	// Write updated JSON data to file
	err = ioutil.WriteFile(jsonPath, newjsonData, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func UnmarshalJson() (error, []map[string]string) {
	jsonData, err := ioutil.ReadFile(jsonPath)

	if err != nil {
		log.Fatal(err)
	}

	var data []map[string]string

	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		log.Fatal(err)
	}
	return err, data
}
