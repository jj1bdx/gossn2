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

type Eisn struct {
	Year  uint
	Month uint
	Date  uint
	SN    int
}

type DailyAp struct {
	Year  uint
	Month uint
	Date  uint
	ApDay int
}

type Flux struct {
	Year    uint
	Month   uint
	Date    uint
	F107Obs float64
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

	// Determine valid daily EISN
	eisnVal := &Eisn{
		SN:    -1,
		Year:  0,
		Month: 0,
		Date:  0,
	}
	for i := numlines - 1; i >= 0; i-- {
		d := datalist[i]
		if d.SN >= 0 {
			eisnVal.SN = d.SN
			eisnVal.Year = d.Year
			eisnVal.Month = d.Month
			eisnVal.Date = d.Date
			break
		}
	}

	// Determine valid daily AP
	dailyApVal := &DailyAp{
		ApDay: -1,
		Year:  0,
		Month: 0,
		Date:  0,
	}
	for i := numlines - 1; i >= 0; i-- {
		d := datalist[i]
		if d.ApDay >= 0 {
			dailyApVal.ApDay = d.ApDay
			dailyApVal.Year = d.Year
			dailyApVal.Month = d.Month
			dailyApVal.Date = d.Date
			break
		}
	}

	// Determine valid daily F10.7 flux value
	fluxVal := &Flux{
		F107Obs: -1.0,
		Year:    0,
		Month:   0,
		Date:    0,
	}
	for i := numlines - 1; i >= 0; i-- {
		d := datalist[i]
		if d.F107Obs >= 0.0 {
			fluxVal.F107Obs = d.F107Obs
			fluxVal.Year = d.Year
			fluxVal.Month = d.Month
			fluxVal.Date = d.Date
			break
		}
	}

	// Pick up the latest valid data
	// If the last line data is not filled,
	// use the one line data before the last line data
	var checkdata Dailydata
	var checkdata2 Dailydata
	if datalist[numlines-1].Ap[0] == -1 {
		checkdata = datalist[numlines-2]
		checkdata2 = datalist[numlines-3]
	} else {
		checkdata = datalist[numlines-1]
		checkdata2 = datalist[numlines-2]
	}

	// For testing purpose: print parsed results
	d := checkdata2
	fmt.Printf("%04d-%02d-%02d "+
		"%.1f/%.1f/%.1f/%.1f/%.1f/%.1f/%.1f/%.1f "+
		"%d/%d/%d/%d/%d/%d/%d/%d "+
		"%d %d %.1f\n",
		d.Year, d.Month, d.Date,
		d.Kp[0], d.Kp[1], d.Kp[2], d.Kp[3],
		d.Kp[4], d.Kp[5], d.Kp[6], d.Kp[7],
		d.Ap[0], d.Ap[1], d.Ap[2], d.Ap[3],
		d.Ap[4], d.Ap[5], d.Ap[6], d.Ap[7],
		d.ApDay, d.SN, d.F107Obs)
	d = checkdata
	fmt.Printf("%04d-%02d-%02d "+
		"%.1f/%.1f/%.1f/%.1f/%.1f/%.1f/%.1f/%.1f "+
		"%d/%d/%d/%d/%d/%d/%d/%d "+
		"%d %d %.1f\n",
		d.Year, d.Month, d.Date,
		d.Kp[0], d.Kp[1], d.Kp[2], d.Kp[3],
		d.Kp[4], d.Kp[5], d.Kp[6], d.Kp[7],
		d.Ap[0], d.Ap[1], d.Ap[2], d.Ap[3],
		d.Ap[4], d.Ap[5], d.Ap[6], d.Ap[7],
		d.ApDay, d.SN, d.F107Obs)

	fmt.Printf("EISN: %04d-%02d-%02d: %d\n",
		eisnVal.Year, eisnVal.Month, eisnVal.Date, eisnVal.SN)
	fmt.Printf("APDay: %04d-%02d-%02d: %d\n",
		dailyApVal.Year, dailyApVal.Month, dailyApVal.Date,
		dailyApVal.ApDay)
	fmt.Printf("Flux: %04d-%02d-%02d: %.1f\n",
		fluxVal.Year, fluxVal.Month, fluxVal.Date,
		fluxVal.F107Obs)

}
