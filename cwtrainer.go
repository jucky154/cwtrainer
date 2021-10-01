/*
 provides the zLog programming interface to the Go language.
 Copyright (C) 2020 JA1ZLO.
*/
package main

import (
	"bufio"
	_ "embed"
	"math"
	"github.com/faiface/beep/speaker"
	"strings"
	"os"
	"time"
)

var (
	calls   [1000]string
	nowcall string
	callall int
	answer  bool
)

//go:embed cwtrainer.dat
var cityMultiList string

func init() {
	CityMultiList = cityMultiList
	OnLaunchEvent = onLaunchEvent
	OnFinishEvent = onFinishEvent
	OnAttachEvent = onAttachEvent
	OnInsertEvent = onInsertEvent
	OnVerifyEvent = onVerifyEvent
	OnPointsEvent = onPointsEvent
}

const (
	rise = 0.1
	fall = 0.1
)

const freq = 800
const rate = 8000
const wpms = 30
const unit = rate * 60.0 / (wpms * 50.0)

//go:embed table.dat
var morse string
var table = make(map[rune]string)

type MorseTone struct {
	done chan bool
	buff [][2]float64
	t int
}

func (tone *MorseTone) Stream(buf [][2]float64) (int, bool) {
	d := len(buf)
	if tone.t + d > len(tone.buff) {
		d = len(tone.buff) - tone.t
	}
	copy(buf[0:d], tone.buff[tone.t:tone.t + d])
	tone.t += d
	if d == 0 {
		tone.done <- true
	}
	return len(buf), d > 0
}

func (tone *MorseTone) Err() error {
	return nil
}

func Play(text string) {
	tone := MorseCodeToMorseTone(RawStringToMorseCode(text))
	err := speaker.Init(rate, 256)
	if err != nil {
		DisplayToast(err.Error())
	}
	speaker.Play(tone)
	<-tone.done
	speaker.Close()
}

func RawStringToMorseCode(text string) string {
	code := []string{}
	for _, ch := range text {
		code = append(code, table[ch])
	}
	return strings.Join(code, " ")
}

func size(code string) (size int) {
	for _, ch := range code {
		switch ch {
		case ' ':
			size += int(unit * 3)
		case ';':
			size += int(unit * 1)
		case '_':
			size += int(unit * 4)
		case '.':
			size += int(unit * 2)
		}
	}
	return size
}

func MorseCodeToMorseTone(code string) *MorseTone {
	idx := 0
	buff := make([][2]float64, size(code))
	for _, ch := range code {
		switch ch {
		case ' ':
			idx = mute(3, idx, buff)
		case ';':
			idx = mute(1, idx, buff)
		case '_':
			idx = tone(3, idx, buff)
			idx = mute(1, idx, buff)
		case '.':
			idx = tone(1, idx, buff)
			idx = mute(1, idx, buff)
		}
	}
	tone := new(MorseTone)
	tone.done = make(chan bool)
	tone.buff = buff
	return tone
}

func tone(time, idx int, buff [][2]float64) int {
	step := 2 * math.Pi * freq / rate
	t1 := int(unit * rise)
	t2 := int(unit * fall)
	t3 := int(unit * time)
	for t := 0; t < t3; t++ {
		amp, r := 1.0, t3-t
		if t < t1 {
			amp *= float64(t) / float64(t1)
		} else if r < t2 {
			amp *= float64(r) / float64(t2)
		}
		buff[idx][0] = math.Sin(float64(t)*step) * amp
		buff[idx][1] = buff[idx][0]
		idx += 1
	}
	return idx
}

func mute(time, idx int, buff [][2]float64) int {
	for t := 0; t < int(unit*time); t++ {
		buff[idx][0] = 0
		buff[idx][1] = buff[idx][0]
		idx += 1
	}
	return idx
}



func onLaunchEvent() {
	reader := strings.NewReader(morse)
	stream := bufio.NewScanner(reader)
	for stream.Scan() {
		val := stream.Text()
		table[rune(val[0])] = val[1:len(val)]
	}
	DisplayToast("CQ!")
}

func onFinishEvent() {
	DisplayToast("Bye")
}

func readtext() {
	//open text file
	fp, err := os.Open("callsigns.txt")
	if err != nil {
		DisplayToast("not find textfile")
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)

	i := 0
	for scanner.Scan() {
		i = i + 1
		//dealing with callsigns one by one
		calls[i] = strings.Trim(scanner.Text(), " ")
	}
	callall = i
}


func checkcw() {
	for i := 1; i < callall+1; i++ {
		nowcall = calls[i]
		answer = false
		for {
			Play(nowcall)
			if answer == true {
				break
			}
			time.Sleep(3 * time.Second)
			if answer == true {
				break
			}
		}
	}
	DisplayToast("finish")
}

func onAttachEvent(test string, path string) {
	readtext()
	go checkcw()
}


func onInsertEvent(qso *QSO) {
	answer = true
}

func onVerifyEvent(qso *QSO)  {
	gscall := qso.GetCall()
	if gscall == nowcall {
		qso.Score = 1
		qso.SetMul1("OK")
		qso.SetNote("OK")
	} else {
		qso.Score = 0
		qso.SetMul1("BAD")
		qso.SetNote("callsign is " + nowcall)
	}
}

func onPointsEvent(score, mults int)int {
	return score
}