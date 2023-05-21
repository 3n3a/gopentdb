package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/3n3a/gopentdb"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main()  {
	o := gopentdb.New(gopentdb.Config{
		BaseUrl: "https://opentdb.com",
	})

	// categories list
	categories, err := o.GetCategories()
	check(err)

	for _, category := range categories {
		questions := make([]gopentdb.Question, 0)

		token, err := o.GetSessionToken()
		check(err)

		count, err := o.GetCategoryCount(category.Id)
		check(err)

		remaining := count
	
		for {
			questionChunk, err := o.GetQuestions(gopentdb.QuestionParams{
				Amount: remaining,
				Category: category.Id,
				Token: token,
			})
			fmt.Printf("Category %d - %s, Chunk Count %d\n", category.Id, category.Name, len(questionChunk))
			check(err)
			questions = append(questions, questionChunk...)

			remaining = remaining - int64(len(questionChunk))

			fmt.Printf("Questions Count %d of %d\n", len(questions), count)
			if int64(len(questions)) == count {
				break
			}
		}

		file, err := json.MarshalIndent(questions, "", " ")
		check(err)
		err = ioutil.WriteFile(fmt.Sprintf("%d-category.json", category.Id), file, 0644)
		check(err)
		fmt.Printf("Saved Questions for %s\n", category.Name)
	}
}