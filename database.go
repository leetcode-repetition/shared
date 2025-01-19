package shared

import (
	"encoding/json"
	"os"

	"log"

	"github.com/supabase-community/supabase-go"
)

func CreateSupabaseClient() (*supabase.Client, error) {
	client, err := supabase.NewClient(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_KEY"), &supabase.ClientOptions{})

	if err != nil {
		log.Printf("Cannot initialize client: %v", err)
	} else {
		log.Printf("Initialized supabase client")
	}

	return client, err
}

func UpsertProblemIntoDatabase(username string, problem LeetCodeProblem) error {
	client, err := CreateSupabaseClient()
	if err != nil {
		return err
	}

	table := os.Getenv("SUPABASE_TABLE")
	_, _, err = client.From(table).
		Upsert(map[string]interface{}{
			"username":           username,
			"titleSlug":          problem.TitleSlug,
			"link":               problem.Link,
			"repeatDate":         problem.RepeatDate,
			"lastCompletionDate": problem.LastCompletionDate,
		}, "username,titleSlug", "", "").
		Execute()

	if err != nil {
		log.Printf("Error upserting database: %v", err)
	}

	log.Printf("Successfully upserted database entry for user: %s", username)
	return err
}

func DeleteProblemFromDatabase(username string, problem_title_slug string) error {
	client, err := CreateSupabaseClient()
	if err != nil {
		return err
	}

	table := os.Getenv("SUPABASE_TABLE")
	_, _, err = client.From(table).
		Delete("", "").
		Eq("username", username).
		Eq("titleSlug", problem_title_slug).
		Execute()

	if err != nil {
		log.Printf("Error deleting database entry: %v", err)
	}

	log.Printf("Successfully deleted database entry for user: %s", username)
	return err
}

func GetProblemsFromDatabase(username string) []LeetCodeProblem {
	var problems []LeetCodeProblem

	client, e := CreateSupabaseClient()
	if e != nil {
		log.Printf("Error creating supabase client: %v", e)
		return []LeetCodeProblem{}
	}
	table := os.Getenv("SUPABASE_TABLE")

	log.Printf("Getting problems from database for user: %s", username)
	rawData, _, err := client.From(table).Select("*", "", false).Eq("username", username).Execute()
	if err != nil {
		log.Printf("Error fetching data: %v", err)
		return []LeetCodeProblem{}
	}

	log.Printf("Raw data: %s", string(rawData))

	var rawProblems []map[string]interface{}
	err = json.Unmarshal(rawData, &rawProblems)
	if err != nil {
		log.Printf("Error unmarshaling data: %v", err)
		return []LeetCodeProblem{}
	}
	for _, rawProblem := range rawProblems {
		problem := LeetCodeProblem{
			Link:               rawProblem["link"].(string),
			TitleSlug:          rawProblem["titleSlug"].(string),
			RepeatDate:         rawProblem["repeatDate"].(string),
			LastCompletionDate: rawProblem["lastCompletionDate"].(string),
		}
		problems = append(problems, problem)
	}

	log.Printf("Problems for user %s: %+v", username, problems)
	return problems
}
