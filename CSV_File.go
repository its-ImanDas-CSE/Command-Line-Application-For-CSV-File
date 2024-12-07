package main

import (
	"encoding/csv" //To read/write CSV files, For working with CSV (Comma Separated Values) files.
	"errors"       // For handling errors,
	"fmt"          // For printing to the console, For printing something on the screen
	"os"           // For working with files and directories, For interacting with the operating system (like reading/writing files).
	"sort"         // For sorting slices, For sorting data.
	"strconv"      // To convert strings to integers
	"strings"      //To work with string manipulation, For handling text and strings.
)

// This defines a custom data Type called CSV
// CSV variable represents each entry in csv file
type CSV struct {
	SiteID                int
	FxiletID              int
	Name                  string
	Criticality           string
	RelevantComputerCount int
}

var filePath string = "fixlets.csv" //  Global filePath

// LoadCSV reads the CSV file and returns a slice, of CSV structs (each element have Custom CSV Struct data type)
func LoadCSV(filePath string) ([]CSV, error) {
	file, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	reader := csv.NewReader(file)    // Creates a CSV reader to read the file.
	records, err := reader.ReadAll() // Reads all rows of the CSV file into a 2D slice

	if err != nil {
		return nil, err
	}

	var entries []CSV // It is a Slice
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
			siteID, err := strconv.Atoi(record[0]) // Convert String into Integer
			if err != nil {
				return nil, err
			}
			// fmt.Println(reflect.TypeOf(siteID))
			fxiletID, err := strconv.Atoi(record[1]) // Convert String into Integer
			if err != nil {
				return nil, err
			}

			relevantComputerCount, err := strconv.Atoi(record[4]) // Convert String into Integer
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

func main() { //Main Function, The programm start from here
	enteries, err := LoadCSV(filePath) //It calls LoadCSV to load the CSV data into the enteries slice.

	if err != nil {
		fmt.Println("Error Loading CSV....!! \n", err)
		return
	}

	for { //This is a menu that asks the user to choose an option (1-6)
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

		case 6: // Saves the current state of enteries back to the CSV file by calling SaveEnteries
			err := SaveEnteries(filePath, enteries)
			if err != nil {
				fmt.Println("Error saving CSV:", err)
			} else {
				fmt.Println("Changes saved.")
			}
			return // return statement immediately exits the main function. This also stops the infinite for loop, and the program terminates.
			break

		default:
			fmt.Println("Invalid choice. Please try again.")

		}

	}
}

func SaveEnteries(filePath string, enteries []CSV) error {
	file, err := os.Create(filePath) // This line attempts to create or open a file at the location specified by filePath.
	// will create a new file or overwrite an existing file with the given path.
	if err != nil {
		return err
	}
	defer file.Close() // will close the file once all the operations are completed.

	writer := csv.NewWriter(file) // Creates a CSV writer to write the file.
	defer writer.Flush()          // This ensures that any buffered data is written to the file before the function exits.
	for _, v := range enteries {
		record := []string{ // new slice
			strconv.Itoa(v.SiteID),   // converts the integer to a string.
			strconv.Itoa(v.FxiletID), // converts the integer to a string.
			v.Name,
			v.Criticality,
			strconv.Itoa(v.RelevantComputerCount), // converts the integer to a string.
		}
		if err := writer.Write(record); err != nil { // This line writes the record (which contains a single row of data) to the CSV file using the csv.Writer.
			return err
		}
	}
	return nil
}

/* cSV *[]CSV: This means cSV is a pointer to a slice of CSV structs. By passing a pointer, we allow the function to modify the original slice.*/
func DeleteEntries(cSV *[]CSV, fxiletID int) error {
	for idx, v := range *cSV {
		if v.FxiletID == fxiletID {
			*cSV = append((*cSV)[:idx], (*cSV)[idx+1:]...) // (*cSV)[:idx]: This gives us all the elements before the element at index idx.
			return nil                                     // (*cSV)[idx+1:]: This gives us all the elements after the element at index idx.
		}
	}
	return errors.New("Entery Not Found!!")
}

// The pointer (*) indicates that we're working with the original slice, not a copy of it. This allows us to modify the slice directly, adding a new entry.
// cSV is a pointer to a slice of CSV structs
func AddEntries(cSV *[]CSV, siteID, fxiletID int, name, criticality string, relevantComputerCount int) {
	*cSV = append(*cSV, CSV{ // It adds a new entry (a new CSV struct) to the slice of CSV entries.
		SiteID:                siteID,
		FxiletID:              fxiletID,
		Name:                  name,
		Criticality:           criticality,
		RelevantComputerCount: relevantComputerCount,
	})
}

// The SortEntries function sorts the enteries slice based on the FxiletID field in ascending order.
func SortEntries(enteries []CSV) { // sort.Slice is a built-in Go function that sorts a slice (like enteries) based on a comparison function you provide.
	// The func(i, j int) bool part defines a comparison function. The comparison function receives two indices, i and j, and should return true if the element at i should come before the element at j, and false otherwise.
	sort.Slice(enteries, func(i, j int) bool { // enteries is the slice that you want to sort.
		return enteries[i].FxiletID < enteries[j].FxiletID
		//It compares the FxiletID field of the two elements at indices i and j in the enteries slice.
		//enteries[i].FxiletID refers to the FxiletID field of the entry at index i.
		//enteries[j].FxiletID refers to the FxiletID field of the entry at index j.
		//The < operator checks if the FxiletID of the entry at i is less than the FxiletID of the entry at j.
		//If it is true, this means the entry at i should come before the entry at j, so the function returns true.
		//If it is false, the function returns false, meaning the entry at i should come after the entry at j.
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
