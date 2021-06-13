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
	"github.com/nextzlog/zylo"
	"os"
	"strings"
	"time"
)

var (
	calls   [100]string
	nowcall string
	callall int
	answer  bool
)

//go:embed cwtrainer.dat
var cwtrainer_list string

func zcities() string {
	return cwtrainer_list
}

func zlaunch() {
	zylo.Notify("CQ!")
}

func zfinish() {
	zylo.Notify("Bye")
}

func readtext() {
	//open text file
	p, _ := os.Getwd()
	fp, err := os.Open(p + "\\wavfiles\\callsigns.txt")
	if err != nil {
		zylo.Notify("panic")
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
		zylo.Notify("panic")
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
	zylo.Notify("finish")
}

func zattach(test string, path string) {
	readtext()
	go checkcw()
}

func zdetach() {
}

func zinsert(qso *zylo.QSO) {
	answer = true
}

func zdelete(qso *zylo.QSO) {
}

func zverify(qso *zylo.QSO) {
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

func zpoints(score, mults int) int {
	return score
}

func zeditor(key int, name string) bool {
	return false
}

func zbutton(btn int, name string) bool {
	return false
}

func main() {}
