package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
)

var urlGFZNowCast string = "https://kp.gfz-potsdam.de/app/files/Kp_ap_Ap_SN_F107_nowcast.txt"

type Dailydata struct {
	Year       uint
	Month      uint
	Date       uint
	UTDays     uint
	UTDaysm    float64
	BSR        uint
	DayBSR     uint
	Kp         [8]float64
	Ap         [8]int
	ApDay      int
	SN         int
	F107Obs    float64
	F107Adj    float64
	Definitive uint
}

const maxDatalines uint = 40
const fieldNumbers int = 28

func main() {
	// Obtain GFZ Nowcast data
	gfzData := getURLBody(urlGFZNowCast)
	gfzScanner := bufio.NewScanner(bytes.NewReader(gfzData[:]))
	// Break the obtained data into lines
	var lines = make([]string, maxDatalines)
	numlines := 0
	for gfzScanner.Scan() {
		line := gfzScanner.Text()
		if err := gfzScanner.Err(); err != nil {
			log.Print(err)
			fmt.Println("Unable to read GFZ Nowcast records")
		} else {
			// Remove comment lines
			if (len(line) > 0) && (line[0:1] != "#") {
				lines[numlines] = line
				numlines++
			}
		}
	}
	// Parse lines
	var datalist = make([]Dailydata, numlines)
	for i := 0; i < numlines; i++ {
		var d Dailydata
		line := lines[i]
		format :=
			"%d %d %d %d %f %d %d " +
				"%f %f %f %f %f %f %f %f " +
				"%d %d %d %d %d %d %d %d " +
				"%d %d %f %f %d"
		n, err := fmt.Sscanf(line, format,
			&d.Year, &d.Month, &d.Date, &d.UTDays, &d.UTDaysm,
			&d.BSR, &d.DayBSR,
			&d.Kp[0], &d.Kp[1], &d.Kp[2], &d.Kp[3],
			&d.Kp[4], &d.Kp[5], &d.Kp[6], &d.Kp[7],
			&d.Ap[0], &d.Ap[1], &d.Ap[2], &d.Ap[3],
			&d.Ap[4], &d.Ap[5], &d.Ap[6], &d.Ap[7],
			&d.ApDay, &d.SN, &d.F107Obs, &d.F107Adj, &d.Definitive)
		if err != nil {
			log.Print(err)
			fmt.Println("Unable to parse GFZ records")
		}
		if n != fieldNumbers {
			log.Print(err)
			fmt.Println("Incorrect field number in GFZ records")
		}
		datalist[i] = d
	}
	// For testing purpose: print parsed results
	for i := 0; i < numlines; i++ {
		d := datalist[i]
		fmt.Printf("%d %04d/%02d/%02d "+
			"%.1f/%.1f/%.1f/%.1f/%.1f/%.1f/%.1f/%.1f "+
			"%d/%d/%d/%d/%d/%d/%d/%d "+
			"%d %d %.1f %.1f\n",
			i,
			d.Year, d.Month, d.Date,
			d.Kp[0], d.Kp[1], d.Kp[2], d.Kp[3],
			d.Kp[4], d.Kp[5], d.Kp[6], d.Kp[7],
			d.Ap[0], d.Ap[1], d.Ap[2], d.Ap[3],
			d.Ap[4], d.Ap[5], d.Ap[6], d.Ap[7],
			d.ApDay, d.SN, d.F107Obs, d.F107Adj)
	}
}
