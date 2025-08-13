package shared

import (
	"encoding/json"
	"fmt"
	"os"

	"log"

	"github.com/supabase-community/supabase-go"
)

var supabaseClient *supabase.Client

func InitSupabaseClient() error {
	var err error
	supabaseClient, err = supabase.NewClient(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_KEY"), &supabase.ClientOptions{})
	if err != nil {
		log.Printf("Cannot initialize client: %v", err)
		return err
	}
	log.Printf("Initialized supabase client")
	return nil
}

func UpsertProblemIntoDatabase(userId string, problem LeetCodeProblem) error {
	if supabaseClient == nil {
		return fmt.Errorf("supabase client not initialized")
	}

	table := os.Getenv("SUPABASE_TABLE")
	_, _, err := supabaseClient.From(table).
		Upsert(map[string]interface{}{
			"userId":             userId,
			"titleSlug":          problem.TitleSlug,
			"link":               problem.Link,
			"repeatDate":         problem.RepeatDate,
			"lastCompletionDate": problem.LastCompletionDate,
		}, "userId,titleSlug", "", "").
		Execute()

	if err != nil {
		log.Printf("Error upserting database: %v", err)
		return err
	}

	log.Printf("Successfully upserted database entry for user: %s", userId)
	return err
}

func UpsertApiKeyIntoDatabase(userId string, token string, apiKey string, apiKeyCreationTime string) error {
	if supabaseClient == nil {
		return fmt.Errorf("supabase client not initialized")
	}

	table := os.Getenv("SUPABASE_API_TABLE")
	_, _, err := supabaseClient.From(table).
		Upsert(map[string]interface{}{
			"userId":             userId,
			"token":              token,
			"apiKey":             apiKey,
			"apiKeyCreationTime": apiKeyCreationTime,
		}, "userId,token", "", "").
		Execute()

	if err != nil {
		log.Printf("Error upserting database: %v", err)
	}

	log.Printf("Successfully upserted database entry for user: %s", userId)
	return err
}

func DeleteProblemFromDatabase(userId string, problem_title_slug string) error {
	if supabaseClient == nil {
		return fmt.Errorf("supabase client not initialized")
	}

	table := os.Getenv("SUPABASE_TABLE")
	_, _, err := supabaseClient.From(table).
		Delete("", "").
		Eq("userId", userId).
		Eq("titleSlug", problem_title_slug).
		Execute()

	if err != nil {
		log.Printf("Error deleting database entry: %v", err)
	}

	log.Printf("Successfully deleted database entry for user: %s", userId)
	return err
}

func GetProblemsFromDatabase(userId string) []LeetCodeProblem {
	if supabaseClient == nil {
		log.Printf("supabase client not initialized")
		return []LeetCodeProblem{}
	}

	var problems []LeetCodeProblem
	table := os.Getenv("SUPABASE_TABLE")

	log.Printf("Getting problems from database for user: %s", userId)
	rawData, _, err := supabaseClient.From(table).Select("*", "", false).Eq("userId", userId).Execute()
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

	log.Printf("Problems for user %s: %+v", userId, problems)
	return problems
}

func GetApiKeyFromDatabase(userId string, token string) (string, string) {
	if supabaseClient == nil {
		log.Printf("supabase client not initialized")
		return "", ""
	}

	table := os.Getenv("SUPABASE_API_TABLE")

	log.Printf("Getting problems from database for user: %s", userId)
	rawData, _, err := supabaseClient.From(table).Select("*", "", false).Eq("userId", userId).Eq("token", token).Execute()
	if err != nil {
		log.Printf("Error fetching data: %v", err)
		return "", ""
	}

	log.Printf("Raw data: %s", string(rawData))

	var results []map[string]interface{}
	err = json.Unmarshal(rawData, &results)
	if err != nil {
		log.Printf("Error unmarshaling data: %v", err)
		return "", ""
	}

	if len(results) == 0 {
		return "", ""
	}

	apiKey, okApiKey := results[0]["apiKey"].(string)
	apiKeyCreationTime, okApiKeyCreationTime := results[0]["apiKeyCreationTime"].(string)
	if !okApiKey || !okApiKeyCreationTime {
		return "", ""
	}

	return apiKey, apiKeyCreationTime
}

func DeleteAllProblemsFromDatabase(userId string) error {
	if supabaseClient == nil {
		return fmt.Errorf("supabase client not initialized")
	}

	table := os.Getenv("SUPABASE_TABLE")
	_, _, err := supabaseClient.From(table).
		Delete("", "").
		Eq("userId", userId).
		Execute()

	if err != nil {
		log.Printf("Error deleting database entry: %v", err)
	}

	log.Printf("Successfully deleted all database entries for user: %s", userId)
	return err
}
