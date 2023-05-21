package main

import (
	"fmt"

	"github.com/3n3a/gopentdb"
)

func main() {
	o := gopentdb.New(gopentdb.Config{
		BaseUrl: "https://opentdb.com",
	})
	fmt.Printf("== Example OpenTdb ==\n==============\n")

	// // List of Categories
	// categories, err := o.GetCategories()
	// if err != nil {
	// 	panic(err)
	// } 
	// fmt.Printf("Categories: %+v\n", categories)


	// // Category Count
	// for _, category := range categories {
	// 	count, err := o.GetCategoryCount(category.Id)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	fmt.Printf("Count for Category %d is %d\n", category.Id, count)
	// }

	// Get Questions for any category
	questions, err := o.GetQuestions(gopentdb.QuestionParams{
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Questions: %+v\n", questions)
}