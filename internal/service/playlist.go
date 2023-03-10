package service

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type Player struct {
	Playlist IPlaylist
}

type IPlaylist interface {
	RunPlayer()
	Play() error
	Pause() error
	AddSong(*bufio.Scanner)
	Next() error
	Prev() error
	CreateSong(*Song)
	LoadSongs()
	WriteSongToJson(*Song) error
	ShowSong(*Song) (map[string]string, error)
	UpdateSong(*Song) *Song
	DeleteSong(*Song) error
}

type Playlist struct {
	Songs    SongList
	mu       *sync.Mutex
	Timer    *Timer
	nextSong chan bool
	prevSong chan bool
	start    chan bool
	IsPlay   bool
}

type Song struct {
	ID       int
	Name     string
	Duration time.Duration
	next     *Song
	prev     *Song
}

type SongList struct {
	Head    *Song
	Current *Song
	tail    *Song
}

type JsonData map[string]string

func NewPlaylist() *Playlist {
	return &Playlist{
		Songs:    SongList{},
		mu:       &sync.Mutex{},
		nextSong: make(chan bool),
		prevSong: make(chan bool),
		start:    make(chan bool),
		IsPlay:   false,
		Timer:    NewTimer(0),
	}
}

func NewPlayer() *Player {
	return &Player{
		Playlist: NewPlaylist(),
	}
}

func (p *Playlist) Play() error {
	if p.IsPlay {
		fmt.Println("Музыка уже играет")
		return errors.New("Музыка уже играет")
	}
	if p.Songs.Head == nil {
		fmt.Println("Плейлист пуст")
		return errors.New("Плейлист пуст")
	}
	p.start <- true
	return nil
}

func (p *Playlist) RunPlayer() {
	for {
		if p.IsPlay {
			p.Timer.Start()
			fmt.Printf("Начала играть песня: %s\n", p.Songs.Current.Name)
		}

		select {
		case <-p.start:
			if p.Songs.Current == nil {
				p.Songs.Current = p.Songs.Head
			}

			p.IsPlay = true
			if p.Timer.timePassed != 0 {
				p.Timer = NewTimer((p.Timer.duration - p.Timer.timePassed).Abs())
			} else {
				p.Timer = NewTimer(p.Songs.Current.Duration.Abs())
			}
		case <-p.Timer.Done():
			fmt.Printf("Песня %s закончилась!\n", p.Songs.Current.Name)
			fmt.Printf("Включается следующая песня!\n")

			if p.Songs.Current.next != nil {
				p.Songs.Current = p.Songs.Current.next
			} else {
				p.Songs.Current = p.Songs.Head
			}

			p.Timer = NewTimer(p.Songs.Current.Duration.Abs())
		case <-p.Timer.Paused():
			p.IsPlay = false

			fmt.Printf("Песня %s была остановлена, осталось время воспроизведения: %s\n",
				p.Songs.Current.Name, p.Timer.duration-p.Timer.timePassed)
		case <-p.nextSong:
			if p.Songs.Current.next != nil {
				p.Songs.Current = p.Songs.Current.next
				p.Timer = NewTimer(p.Songs.Current.Duration.Abs())
			} else {
				p.Timer = NewTimer(p.Songs.Current.Duration.Abs())
			}
			if !p.IsPlay {
				p.IsPlay = true
			}

		case <-p.prevSong:
			if p.Songs.Current.prev != nil {
				p.Songs.Current = p.Songs.Current.prev
				p.Timer = NewTimer(p.Songs.Current.Duration.Abs())
			} else {
				p.Timer = NewTimer(p.Songs.Current.Duration.Abs())
			}
			if !p.IsPlay {
				p.IsPlay = true
			}
		}
	}
}

func (p *Playlist) Pause() error {
	if !p.IsPlay {
		fmt.Println("Воспроизведение плейлиста уже приостановлено!")
		return errors.New("Воспроизведение плейлиста уже приостановлено!")
	}

	p.Timer.Pause()
	return nil
}

func (p *Playlist) AddSong(scan *bufio.Scanner) {
	fmt.Println("Введите название: ")
	scan.Scan()
	Name := scan.Text()

	fmt.Println("Введите длительность: ")
	scan.Scan()
	Duration, err := strconv.Atoi(scan.Text())
	for err != nil {
		fmt.Println("Введите число")
		scan.Scan()
		Duration, err = strconv.Atoi(scan.Text())
	}

	p.mu.Lock()

	newSong := &Song{Name: Name, Duration: time.Duration(Duration).Abs() * time.Second, ID: rand.Int()}
	p.CreateSong(newSong)

	newSongJson := &Song{Name: Name, Duration: time.Duration(Duration).Abs(), ID: rand.Int()}
	err = p.WriteSongToJson(newSongJson)

	if err != nil {
		panic(err)
	}
	p.mu.Unlock()
}

func (p *Playlist) Next() error {
	if p.Songs.Current == nil {
		fmt.Println("Плейлист пуст")
		return errors.New("Плейлист пуст")
	}
	p.nextSong <- true
	return nil
}

func (p *Playlist) Prev() error {
	if p.Songs.Current == nil {
		fmt.Println("Плейлист пуст")
		return errors.New("Плейлист пуст")
	}
	p.prevSong <- true
	return nil
}
