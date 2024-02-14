package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/xuri/excelize/v2"
)

var days = []string{"MONDAY", "TUESDAY", "WEDNESDAY", "THURSDAY", "FRIDAY", "SATURDAY", "SUNDAY"}
var meals = []string{"BREAKFAST", "LUNCH", "DINNER"}

func isDay(s string) bool {
	for _, day := range days {
		if s == day {
			return true
		}
	}
	return false
}

func getmealitems(f *excelize.File, targetDay, targetMeal string) []string {
	data, _ := f.GetRows("Sheet1")
	targetDay = strings.ToUpper(targetDay)
	targetMeal = strings.ToUpper(targetMeal)
	var items []string
	var dayindex int
	var mealindexstart int
	dayfound := false
	mealfound := false
	var maxdown int
	var maxright int = len(data[0]) - 1
	for index, row := range data {
		if len(row) < maxright {
			maxright = len(row) - 1
			maxdown = index
			break
		}
	}

	firstRow := data[0]
	for i, cell := range firstRow {
		if cell == targetDay {
			dayindex = i
			dayfound = true
			break
		}
	}
	if !dayfound {
		fmt.Println("Day not found")
		return items
	}

	for i := 0; i < len(data); i++ {
		if data[i][dayindex] == targetMeal {
			mealindexstart = i + 1
			mealfound = true
			break
		}
	}
	if !mealfound {
		fmt.Println("Meal not found")
		return items
	}

	for j := mealindexstart; j < len(data) && data[j][dayindex] != "" && !isDay(data[j][dayindex]); j++ {
		items = append(items, data[j][dayindex])
		if j == maxdown-1 && dayindex > maxright {
			break
		}
	}
	return items

}

func getmealitemcount(f *excelize.File, targetDay, targetMeal string) int {
	return len(getmealitems(f, targetDay, targetMeal))
}

func itemchecker(f *excelize.File, targetDay, targetMeal, item string) bool {
	item = strings.ToUpper(item)
	for _, searched := range getmealitems(f, targetDay, targetMeal) {
		if item == strings.ToUpper(searched) {
			return true
		}
	}
	return false
}

type Menu struct {
	Day   string
	Date  string
	Meal  string
	Items []string
}

func makestruct(f *excelize.File) []Menu {
	data, _ := f.GetRows("Sheet1")
	var Mess []Menu
	dayrow := data[0]
	daterow := data[1]
	//create a new menu object with the day and date and meal and items by calling getmealitems
	//append this menu object to the Mess array
	for i := 0; i < len(dayrow); i++ {
		if dayrow[i] != "" {
			day := dayrow[i]
			date := daterow[i]
			for _, meal := range meals {
				items := getmealitems(f, day, meal)
				Mess = append(Mess, Menu{Day: day, Date: date, Meal: meal, Items: items})
			}
		}
	}
	return Mess
}

// Create a function that systematically converts the entire menu into json and saves this data as a json file in the same directory
func makejson(f *excelize.File) {
	Mess := makestruct(f)
	//convert the Mess array to json and save it as a json file
	jsonData, _ := json.Marshal(Mess)
	jsonFile, _ := os.Create("menu.json")
	jsonFile.Write(jsonData)
	jsonFile.Close()
	fmt.Print("JSON file created")
}

/*
	Using the above data, create a struct that contains the day, date, meal and the items in that meal.

Furthermore, create a method that prints the details of each meal instance
*/
func (m Menu) printdetails() {
	fmt.Printf("Day: %s, Date: %s, Meal: %s, Items: %v\n", m.Day, m.Date, m.Meal, m.Items)
}

func main() {
	file, _ := excelize.OpenFile("Sample-Menu.xlsx")
	//Menu Driver function that calls the above functions
	var ioflag int = 1
	for ioflag != 0 {
		fmt.Println("1. Get Meal Items")
		fmt.Println("2. Get Meal Item Count")
		fmt.Println("3. Item Checker")
		fmt.Println("4. Make Struct and Print All Meal Instances")
		fmt.Println("5. Make JSON")
		fmt.Println("6. Exit")
		fmt.Println("Enter your choice: ")

		var targetDay, targetMeal, item string

		fmt.Scanln(&ioflag)
		switch ioflag {
		case 1:
			fmt.Println("Enter the day: ")
			fmt.Scanln(&targetDay)
			fmt.Println("Enter the meal: ")
			fmt.Scanln(&targetMeal)
			fmt.Println(getmealitems(file, targetDay, targetMeal))
		case 2:
			fmt.Println("Enter the day: ")
			fmt.Scanln(&targetDay)
			fmt.Println("Enter the meal: ")
			fmt.Scanln(&targetMeal)
			fmt.Println(getmealitemcount(file, targetDay, targetMeal))
		case 3:
			fmt.Println("Enter the day: ")
			fmt.Scanln(&targetDay)
			fmt.Println("Enter the meal: ")
			fmt.Scanln(&targetMeal)
			fmt.Println("Enter the item: ")
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				item = scanner.Text()
			}
			ans := itemchecker(file, targetDay, targetMeal, item)
			if ans {
				fmt.Println("Item found")
			} else {
				fmt.Println("Item not found")
			}
		case 4:
			Structure := makestruct(file)
			for i := 0; i < len(Structure); i++ {
				Structure[i].printdetails()
			}
		case 5:
			makejson(file)
		case 6:
			ioflag = 0
		default:
			fmt.Println("Invalid choice")
		}
	}
}
