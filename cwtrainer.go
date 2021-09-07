/*
 provides the zLog programming interface to the Go language.
 Copyright (C) 2020 JA1ZLO.
*/
package main

import (
	"bufio"
	_ "embed"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"os"
	"strings"
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


func onLaunchEvent() {
	DisplayToast("CQ!")
}

func onFinishEvent() {
	DisplayToast("Bye")
}

func readtext() {
	//open text file
	p, _ := os.Getwd()
	fp, err := os.Open(p + "\\wavfiles\\callsigns.txt")
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

func sound() {
	p, _ := os.Getwd()
	f, err := os.Open(p + "\\wavfiles\\" + nowcall + ".wav")
	if err != nil {
		DisplayToast("not find wavfile")
	}
	st, format, _ := wav.Decode(f)
	defer st.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	done := make(chan bool)
	speaker.Play(beep.Seq(st, beep.Callback(func() {
		done <- true
	})))
	<-done
}

func checkcw() {
	for i := 1; i < callall+1; i++ {
		nowcall = calls[i]
		answer = false
		for {
			sound()
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