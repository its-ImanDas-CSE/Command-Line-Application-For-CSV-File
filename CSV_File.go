package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// struct variable represents each entry in csv file
type CSV struct {
	SiteID                int
	FxiletID              int
	Name                  string
	Criticality           string
	RelevantComputerCount int
}

var filePath string = "fixlets.csv"

// LoadCSV reads the CSV file and returns a slice of Entry structs
func LoadCSV(filePath string) ([]CSV, error) {
	file, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()

	if err != nil {
		return nil, err
	}

	var entries []CSV
	for i, record := range records {

		if i > 0 { //to skip header
			if len(record) < 4 {
				return nil, fmt.Errorf("Invalid data on line %d", i+1)
			}

			// fmt.Println(record[0])
			// fmt.Println(record[1])
			// fmt.Println(record[2])
			// fmt.Println(record[3])
			// fmt.Println(record[4])

			// fmt.Println(reflect.TypeOf(record[0]))
			siteID, err := strconv.Atoi(record[0])
			if err != nil {
				return nil, err
			}
			// fmt.Println(reflect.TypeOf(siteID))
			fxiletID, err := strconv.Atoi(record[1])
			if err != nil {
				return nil, err
			}

			relevantComputerCount, err := strconv.Atoi(record[4])
			if err != nil {
				return nil, err
			}
			entries = append(
				entries, CSV{
					SiteID:                siteID,
					FxiletID:              fxiletID,
					Name:                  record[2],
					Criticality:           record[3],
					RelevantComputerCount: relevantComputerCount,
				})
			// fmt.Println(i)
			// fmt.Println(record)
		}
	}
	// fmt.Println(records)

	return entries, nil
}

func main() {
	enteries, err := LoadCSV(filePath)

	if err != nil {
		fmt.Println("Error Loading CSV....!! \n", err)
		return
	}

	for {
		fmt.Println("\nChoose an option:")
		fmt.Println("1. List entries")
		fmt.Println("2. Query entries")
		fmt.Println("3. Sort entries by fxiletID")
		fmt.Println("4. Add entry")
		fmt.Println("5. Delete entry")
		fmt.Println("6. Exit")

		var choice int

		fmt.Scanln(&choice)

		switch choice {
		case 1:
			ListEnteries(enteries)
			break

		case 2:
			fmt.Print("Enter FxiletID to search: ")
			var query int
			fmt.Scanln(&query)
			result := QueryEnteries(enteries, query)
			if len(result) < 1 {
				fmt.Println("No Data Found for enterd FxiletID!!")
			} else {
				ListEnteries(result)
			}
			break

		case 3:
			SortEntries(enteries)
			fmt.Println("Entries sorted by FxiletID.")
			break

		case 4:
			var siteID, fxiletID, relevantComputerCount int
			var name, criticality string
			fmt.Println("Enter SiteID")
			fmt.Scanln(&siteID)
			fmt.Println("Enter FxiletID")
			fmt.Scanln(&fxiletID)
			fmt.Println("Enter Name")
			fmt.Scanln(&name)
			fmt.Println("Enter Criticality")
			fmt.Scanln(&criticality)
			fmt.Println("Enter RelevantComputerCount")
			fmt.Scanln(&relevantComputerCount)
			AddEntries(&enteries, siteID, fxiletID, name, criticality, relevantComputerCount)
			fmt.Println("Entry added.")
			break

		case 5:
			var fxiletID int
			fmt.Println("Enter FxiletID")
			fmt.Scanln(&fxiletID)
			err = DeleteEntries(&enteries, fxiletID)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Entry deleted.")
			}
			break

		case 6:
			err := SaveEnteries(filePath, enteries)
			if err != nil {
				fmt.Println("Error saving CSV:", err)
			} else {
				fmt.Println("Changes saved.")
			}
			return
			break

		default:
			fmt.Println("Invalid choice. Please try again.")

		}

	}
}

func SaveEnteries(filePath string, enteries []CSV) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	for _, v := range enteries {
		record := []string{
			strconv.Itoa(v.SiteID),
			strconv.Itoa(v.FxiletID),
			v.Name,
			v.Criticality,
			strconv.Itoa(v.RelevantComputerCount),
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	return nil
}

func DeleteEntries(cSV *[]CSV, fxiletID int) error {
	for idx, v := range *cSV {
		if v.FxiletID == fxiletID {
			*cSV = append((*cSV)[:idx], (*cSV)[idx+1:]...)
			return nil
		}
	}
	return errors.New("Entery Not Found!!")
}

func AddEntries(cSV *[]CSV, siteID, fxiletID int, name, criticality string, relevantComputerCount int) {
	*cSV = append(*cSV, CSV{
		SiteID:                siteID,
		FxiletID:              fxiletID,
		Name:                  name,
		Criticality:           criticality,
		RelevantComputerCount: relevantComputerCount,
	})
}

func SortEntries(enteries []CSV) {
	sort.Slice(enteries, func(i, j int) bool {
		return enteries[i].FxiletID < enteries[j].FxiletID
	})
}

func QueryEnteries(enteries []CSV, query int) []CSV {
	var result []CSV
	for _, v := range enteries {
		// if strings.Contains(strings.ToLower(v.Name),strings.ToLower(query))
		if v.FxiletID == query {
			result = append(result, v)
		}
	}
	return result
}

func ListEnteries(enteries []CSV) {
	fmt.Println("SiteID     |     FxiletID    |     Name    |    Criticality     |     RelevantComputerCount")
	fmt.Println(strings.Repeat("-", 40))
	for _, v := range enteries {
		fmt.Printf("%d     |     %d    |     %s    |    %s     |     %d", v.SiteID, v.FxiletID, v.Name, v.Criticality, v.RelevantComputerCount)
	}
}
