package main

import (
	"fmt"
	"os"
)

func main() {

	for _, project := range projects {

		// Create areays for the different types of issues
		var stories []*JiraStory
		var epics []*JiraEpic
		var features []*JiraFeature
		var outcomes []*JiraOutcome
		var tasks []*JiraTask
		var subtasks []*JiraSubTask
		var bugs []string

		// Get all the issues
		allIssues, err := getIssues(project)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for _, issue := range allIssues.Issues {
			kind := issue.Fields.Issuetype.Name

			switch kind {
			case "Bug":
				if !contains(bugs, issue.Key) {
					bugs = append(bugs, issue.Key)
				}
			case "Task":
				if !containsTask(tasks, issue.Key) {
					tasks = append(tasks, &JiraTask{Name: issue.Key})
				}
			case "Sub-task":
				if !containsSubTask(subtasks, issue.Key) {
					subtasks = append(subtasks, &JiraSubTask{Name: issue.Key})
				}
			case "Story":
				if !containsStory(stories, issue.Key) {
					stories = append(stories, &JiraStory{Name: issue.Key})
				}
			case "Epic":
				if !containsEpic(epics, issue.Key) {
					epics = append(epics, &JiraEpic{Name: issue.Key})
				}
			case "Feature":
				if !containsFeature(features, issue.Key) {
					features = append(features, &JiraFeature{Name: issue.Key})
				}
			case "Outcome":
				if !containsOutcome(outcomes, issue.Key) {
					outcomes = append(outcomes, &JiraOutcome{Name: issue.Key})
				}
			}
		}

		// Explicitely for every card
		// Start with the features, in order to make sure we have the Outcomes
		// ------------------------------------------------------------------- //
		for _, issue := range allIssues.Issues {
			kind := issue.Fields.Issuetype.Name

			switch kind {
			case "Feature":
				// Get the Outcome
				outcome := issue.Fields.Outcome
				if outcome == "" {
					fmt.Printf("WARNING: Feature %s is not part of any Outcome!\n", issue.Key)
					continue
				}
				// Check if the Outcome is already in the array of Outcomes
				var outcomeExists bool
				for _, o := range outcomes {
					if o.Name == outcome {
						outcomeExists = true
					}
				}
				// If the Outcome is not in the array of Outcomes, add it
				if !outcomeExists {
					fmt.Printf("Feature %s is part of the Outcome %s that is not related to SANDBOX!\n", issue.Key, outcome)
					outcomes = append(outcomes, &JiraOutcome{Name: outcome})
				}
				// Add the Feature to the Outcome
				for _, o := range outcomes {
					if o.Name == outcome {
						o.Features = append(o.Features, &JiraFeature{Name: issue.Key})
					}
				}
			}
		}

		// ------------------------------------------------------------------- //
		// Continue with epics, to make sure we have the features
		var rerunFeatureLoop bool
		for _, issue := range allIssues.Issues {
			kind := issue.Fields.Issuetype.Name

			switch kind {
			case "Epic":
				// Get the Feature
				feature := issue.Fields.Feature
				if feature.Key == "" {
					fmt.Printf("WARNING: Epic %s is not part of any Feature!\n", issue.Key)
					continue
				}
				// Check if the Feature is already in the array of Features
				var featureExists bool
				for _, f := range features {
					if f.Name == feature.Key {
						featureExists = true
					}
				}
				// If the Feature is not in the array of Features, add it
				if !featureExists {
					fmt.Printf("Epic %s is part of the Feature %s that is not related to SANDBOX!\n", issue.Key, feature.Key)
					features = append(features, &JiraFeature{Name: feature.Key})
					rerunFeatureLoop = true
				}
				// Add the Epic to the Feature
				for _, f := range features {
					if f.Name == feature.Key {
						f.Epics = append(f.Epics, &JiraEpic{Name: issue.Key})
					}
				}
			}
		}

		if rerunFeatureLoop {
			for _, issue := range allIssues.Issues {
				kind := issue.Fields.Issuetype.Name

				switch kind {
				case "Feature":
					// Get the Outcome
					outcome := issue.Fields.Outcome
					if outcome == "" {
						fmt.Printf("WARNING: Feature %s is not part of any Outcome!\n", issue.Key)
						continue
					}
					// Check if the Outcome is already in the array of Outcomes
					var outcomeExists bool
					for _, o := range outcomes {
						if o.Name == outcome {
							outcomeExists = true
						}
					}
					// If the Outcome is not in the array of Outcomes, add it
					if !outcomeExists {
						fmt.Printf("Feature %s is part of the Outcome %s that is not related to SANDBOX!\n", issue.Key, outcome)
						outcomes = append(outcomes, &JiraOutcome{Name: outcome})
					}
					// Add the Feature to the Outcome
					for _, o := range outcomes {
						if o.Name == outcome {
							o.Features = append(o.Features, &JiraFeature{Name: issue.Key})
						}
					}
				}
			}
		}

		// Update the outcomes.features with the features.epics
		for _, f := range features {
			for _, o := range outcomes {
				for _, of := range o.Features {
					if of.Name == f.Name {
						of.Epics = f.Epics
					}
				}
			}
		}

		// ------------------------------------------------------------------- //
		// Continue with stories, to make sure we have the epics
		var rerunEpicLoop bool
		for _, issue := range allIssues.Issues {
			kind := issue.Fields.Issuetype.Name

			switch kind {
			case "Story":
				// Get the Epic
				epic := issue.Fields.Epic
				if epic == "" {
					fmt.Printf("WARNING: Story %s is not part of any Epic!\n", issue.Key)
					continue
				}
				// Check if the Epic is already in the array of Epics
				var epicExists bool
				for _, e := range epics {
					if e.Name == epic {
						epicExists = true
					}
				}
				// If the Epic is not in the array of Epics, add it
				if !epicExists {
					fmt.Printf("Story %s is part of the Epic %s that is not related to SANDBOX!\n", issue.Key, epic)
					epics = append(epics, &JiraEpic{Name: epic})
					rerunEpicLoop = true
				}
				// Add the Story to the Epic
				for _, e := range epics {
					if e.Name == epic {
						e.Stories = append(e.Stories, &JiraStory{Name: issue.Key})
					}
				}
			}
		}

		if rerunEpicLoop {
			for _, issue := range allIssues.Issues {
				kind := issue.Fields.Issuetype.Name

				switch kind {
				case "Epic":
					// Get the Feature
					feature := issue.Fields.Feature
					if feature.Key == "" {
						fmt.Printf("WARNING: Epic %s is not part of any Feature!\n", issue.Key)
						continue
					}
					// Check if the Feature is already in the array of Features
					var featureExists bool
					for _, f := range features {
						if f.Name == feature.Key {
							featureExists = true
						}
					}
					// If the Feature is not in the array of Features, add it
					if !featureExists {
						fmt.Printf("Epic %s is part of the Feature %s that is not related to SANDBOX!\n", issue.Key, feature.Key)
						features = append(features, &JiraFeature{Name: feature.Key})
					}
					// Add the Epic to the Feature
					for _, f := range features {
						if f.Name == feature.Key {
							f.Epics = append(f.Epics, &JiraEpic{Name: issue.Key})
						}
					}
				}
			}
		}

		// Update the features.epics with the epics.stories
		for _, e := range epics {
			for _, o := range outcomes {
				for _, of := range o.Features {
					for _, ef := range of.Epics {
						if ef.Name == e.Name {
							ef.Stories = e.Stories
						}
					}
				}
			}
		}

		// ------------------------------------------------------------------- //

		// Remove duplicates from all the arrays
		var featuresWithoutDuplicates []*JiraFeature
		for _, f := range features {
			if !containsFeature(featuresWithoutDuplicates, f.Name) {
				featuresWithoutDuplicates = append(featuresWithoutDuplicates, f)
			}
		}

		// replace features with featuresWithoutDuplicates
		features = featuresWithoutDuplicates

		var epicsWithoutDuplicates []*JiraEpic
		for _, e := range epics {
			if !containsEpic(epicsWithoutDuplicates, e.Name) {
				epicsWithoutDuplicates = append(epicsWithoutDuplicates, e)
			}
		}

		// replace epics with epicsWithoutDuplicates
		epics = epicsWithoutDuplicates

		var storiesWithoutDuplicates []*JiraStory
		for _, s := range stories {
			if !containsStory(storiesWithoutDuplicates, s.Name) {
				storiesWithoutDuplicates = append(storiesWithoutDuplicates, s)
			}
		}

		// replace stories with storiesWithoutDuplicates
		stories = storiesWithoutDuplicates

		var outcomesWithoutDuplicates []*JiraOutcome
		for _, o := range outcomes {
			if !containsOutcome(outcomesWithoutDuplicates, o.Name) {
				outcomesWithoutDuplicates = append(outcomesWithoutDuplicates, o)
			}
		}

		// replace outcomes with outcomesWithoutDuplicates
		outcomes = outcomesWithoutDuplicates

		// Print how many outcomes, how many feautures, how many epics, how many stories
		fmt.Printf("There are %d outcomes, %d features, %d epics, %d stories, %d tasks, %d sub-tasks, %d bugs\n", len(outcomes), len(features), len(epics), len(stories), len(tasks), len(subtasks), len(bugs))

		// Add all of them (sum) to the total
		total := len(outcomes) + len(features) + len(epics) + len(stories) + len(tasks) + len(subtasks) + len(bugs)
		fmt.Println("Total:", total)

		// Print a tree of the issues, starting from the Outcomes as top-level
		// Continue with the features as second-level.

		// Remove duplicates from Features of Outcomes
		for _, o := range outcomes {
			var featuresWithoutDuplicates []*JiraFeature
			for _, f := range o.Features {
				if !containsFeature(featuresWithoutDuplicates, f.Name) {
					featuresWithoutDuplicates = append(featuresWithoutDuplicates, f)
				}
			}
			o.Features = featuresWithoutDuplicates
		}

		// Remove duplicates from Epics of Features of Outcomes
		for _, o := range outcomes {
			for _, f := range o.Features {
				var epicsWithoutDuplicates []*JiraEpic
				for _, e := range f.Epics {
					if !containsEpic(epicsWithoutDuplicates, e.Name) {
						epicsWithoutDuplicates = append(epicsWithoutDuplicates, e)
					}
				}
				f.Epics = epicsWithoutDuplicates
			}
		}

		// Remove duplicates from Stories of Epics of Features of Outcomes
		for _, o := range outcomes {
			for _, f := range o.Features {
				for _, e := range f.Epics {
					var storiesWithoutDuplicates []*JiraStory
					for _, s := range e.Stories {
						if !containsStory(storiesWithoutDuplicates, s.Name) {
							storiesWithoutDuplicates = append(storiesWithoutDuplicates, s)
						}
					}
					e.Stories = storiesWithoutDuplicates
				}
			}
		}

		// Print the Outcomes
		fmt.Println("Sandbox End-To-End Traceability with VSEMs")
		for _, o := range outcomes {
			fmt.Println("| Outcome: ", o.Name)
			// Print the Features
			for _, f := range o.Features {
				fmt.Println("  | Feature: ", f.Name)
				// Print the Epics
				for _, e := range f.Epics {
					fmt.Println("    | Epic: ", e.Name)
					// Print the Stories
					for _, s := range e.Stories {
						fmt.Println("      | Story: ", s.Name)
						// 		// Print the Tasks
						// 		for _, t := range tasks {
						// 			fmt.Println("        ", t.Name)
						// 			// Print the SubTasks
						// 			for _, st := range subtasks {
						// 				fmt.Println("          ", st.Name)
						// 			}
						// 		}
					}
				}
			}
			fmt.Println()
		}
		fmt.Println("####################")
		fmt.Println("# Untrackable Work #")
		fmt.Println("####################")
		fmt.Println()
		outcomesWithoutFeatures := getOutcomesWithoutFeatures(outcomes)
		if len(outcomesWithoutFeatures) > 0 {
			fmt.Printf("\nOutcomes without features (%v): %v\n", len(outcomesWithoutFeatures), outcomesWithoutFeatures)
		}

		featuresWithoutOutcome := getFeaturesWithoutOutcome(features, outcomes)
		if len(featuresWithoutOutcome) > 0 {
			fmt.Printf("\nFeatures without outcomes (%v): %v\n", len(featuresWithoutOutcome), featuresWithoutOutcome)
		}

		featuresWithoutEpics := getFeaturesWithoutEpics(features)
		if len(featuresWithoutEpics) > 0 {
			fmt.Printf("\nFeatures without epics (%v): %v\n", len(featuresWithoutEpics), featuresWithoutEpics)
		}

		epicsWithoutFeatures := getEpicsWithoutFeatures(epics, features)
		if len(epicsWithoutFeatures) > 0 {
			fmt.Printf("\nEpics without features (%v): %v\n", len(epicsWithoutFeatures), epicsWithoutFeatures)
		}

		epicWithoutStories := getEpicsWithoutStories(epics)
		if len(epicWithoutStories) > 0 {
			fmt.Printf("\nEpics without stories (%v): %v\n", len(epicWithoutStories), epicWithoutStories)
		}

		storiesWithoutEpics := getStoriesWithoutEpics(stories, epics)
		if len(storiesWithoutEpics) > 0 {
			fmt.Printf("\nStories without epics (%v): %v\n", len(storiesWithoutEpics), storiesWithoutEpics)
		}
		fmt.Println()

		// Count all the Outcomes, Features, Epics, Stories, Tasks, SubTasks, Bugs
		totalWork := len(outcomes) + len(features) + len(epics) + len(stories) + len(tasks) + len(subtasks) + len(bugs)

		// Print the total work that is not tracked
		// Put all the untracked work in a slice and get rid of the duplicates
		var untrackedWork []string
		untrackedWork = append(untrackedWork, outcomesWithoutFeatures...)
		untrackedWork = append(untrackedWork, featuresWithoutOutcome...)
		untrackedWork = append(untrackedWork, featuresWithoutEpics...)
		untrackedWork = append(untrackedWork, epicsWithoutFeatures...)
		untrackedWork = append(untrackedWork, epicWithoutStories...)
		untrackedWork = append(untrackedWork, storiesWithoutEpics...)
		// Remove duplicates
		var untrackedWorkWithoutDuplicates []string
		for _, w := range untrackedWork {
			if !contains(untrackedWorkWithoutDuplicates, w) {
				untrackedWorkWithoutDuplicates = append(untrackedWorkWithoutDuplicates, w)
			}
		}

		fmt.Println(" ----------------------------------------------------- ")
		unhealthyPercentage := float32(100*len(untrackedWorkWithoutDuplicates)) / float32(totalWork)
		fmt.Println("Untracked work percentage:", unhealthyPercentage, "%. Total untracked JIRA issues", len(untrackedWorkWithoutDuplicates), "out of", totalWork, "total JIRA issues.")
		// Print Healthy percentage
		healthyPercentage := float32(100) - unhealthyPercentage
		fmt.Println("Tracked work percentage:", healthyPercentage, "%")
		fmt.Println("")
		// Depending on the healthyPercentage, give a grade from A to F
		if healthyPercentage >= 90 {
			fmt.Println("Grade: A")
		} else if healthyPercentage >= 80 {
			fmt.Println("Grade: B")
		} else if healthyPercentage >= 70 {
			fmt.Println("Grade: C")
		} else if healthyPercentage >= 60 {
			fmt.Println("Grade: D")
		} else if healthyPercentage >= 50 {
			fmt.Println("Grade: E")
		} else {
			fmt.Println("Grade: F")
		}

		// Check if story 'SANDBOX-368' is part of any epic, and if it is check if it is part of any feature and if it is check if it's part of any outcome
		// If it is not part of any epic, feature or outcome, print it

		allUpdatedStories, err := getLast7dUpdates(project)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		listofUpdatedStories := []string{}

		// get list of Stories from for _, issue := range allIssues.Issues {
		for _, updatedStory := range allUpdatedStories.Issues {
			kind := updatedStory.Fields.Issuetype.Name
			switch kind {
			case "Story":
				listofUpdatedStories = append(listofUpdatedStories, updatedStory.Key)
			}
		}

		fmt.Println()
		fmt.Println("In the last 7 days, ", len(listofUpdatedStories), "stories have been updated:", listofUpdatedStories)
		fmt.Println("-----------------------------------------------------------------------------")

		// create a map called reportOutcomes that will have the outcome as key and the list of stories as value
		reportOutcomes := make(map[string][]string)

		// create reportedOutcomed array to keep track of the reported outcomes
		var reportedOutcomes []string

		for _, story := range listofUpdatedStories {
			fmt.Printf("Checking health of story: %v", story)

			isStoryPartOfOutcome, storysOutcome, err := checkStoryHealth(story, epics, features, outcomes)
			if isStoryPartOfOutcome {
				fmt.Printf(" --> [OK] : story %s is part of outcome %s\n", story, storysOutcome)
				// Update the reportOurcomes map
				reportOutcomes[storysOutcome] = append(reportOutcomes[storysOutcome], story)
				// Append the reportedOutcome to the reportedOutcomes array
				reportedOutcomes = append(reportedOutcomes, storysOutcome)
			} else {
				fmt.Printf(" --> [BAD]: %v\n", err)
			}
		}

		// Report
		// Loop through reportedOutcomes and match them with the reportOutcomes map
		fmt.Println()
		fmt.Println("Report: ", project)
		fmt.Println("-------")
		for _, reportedOutcome := range reportedOutcomes {
			fmt.Println()
			if len(reportOutcomes[reportedOutcome]) != 0 {
				fmt.Println("* Outcome:", reportedOutcome)
				for _, v := range reportOutcomes[reportedOutcome] {
					s, err := getStoryInfo(v)
					if err != nil {
						fmt.Println(err)
					}
					title := s.Fields.Summary
					status := s.Fields.Status.Name
					fmt.Printf("  - Story: %s\t[%s]\t-\t%s\n", v, status, title)
				}
			}

		}
	}

}

func checkStoryHealth(story string, epics []*JiraEpic, features []*JiraFeature, outcomes []*JiraOutcome) (bool, string, error) {
	var storyIsPartOfAnyEpic bool
	var storyIsPartOfAnyFeature bool
	var storyIsPartOfAnyOutcome bool
	var err error
	var ep, feat, out string

	for _, e := range epics {
		for _, s := range e.Stories {
			if s.Name == story {
				storyIsPartOfAnyEpic = true
				ep = e.Name
			}
		}
	}
	if !storyIsPartOfAnyEpic {
		err = fmt.Errorf("story %s is not part of any epic", story)
		return false, "", err
	}

	for _, f := range features {
		for _, s := range f.Epics {
			for _, st := range s.Stories {
				if st.Name == story {
					storyIsPartOfAnyFeature = true
					feat = f.Name
				}
			}
		}
	}
	if !storyIsPartOfAnyFeature {
		err = fmt.Errorf("story %s is part of the epic %v which doesn't belong to any feature", story, ep)
		return false, "", err
	}
	for _, o := range outcomes {
		for _, f := range o.Features {
			for _, s := range f.Epics {
				for _, st := range s.Stories {
					if st.Name == story {
						storyIsPartOfAnyOutcome = true
						out = o.Name
					}
				}
			}
		}
	}

	if !storyIsPartOfAnyOutcome {
		err = fmt.Errorf("story %s is part of epic %v which belong to %v feature, which does not belong any outcome", story, ep, feat)
		return false, "", err
	}

	return true, out, nil
}

// printEpicsWithoutStories prints the epics without stories
func getEpicsWithoutStories(epics []*JiraEpic) []string {
	listOfEpicsWithoutStories := []string{}
	// Print the Epics
	for _, e := range epics {
		if len(e.Stories) == 0 {
			listOfEpicsWithoutStories = append(listOfEpicsWithoutStories, e.Name)
		}
	}

	return listOfEpicsWithoutStories
}

// printFeaturesWithoutOutcome prints the features without outcomes
func getFeaturesWithoutOutcome(features []*JiraFeature, outcomes []*JiraOutcome) []string {
	listOfFeaturesWithoutOutcome := []string{}
	listOfFeaturesWithOutcomes := []string{}
	for _, f := range features {
		for _, o := range outcomes {
			for _, of := range o.Features {
				if of.Name == f.Name {
					listOfFeaturesWithOutcomes = append(listOfFeaturesWithOutcomes, f.Name)
				}
			}
		}
	}

	for _, f := range features {
		// If the f.Name is not part of the listOfFeaturesWithOutcomes array, add it to the listOfFeaturesWithoutOutcome
		if !contains(listOfFeaturesWithOutcomes, f.Name) {
			listOfFeaturesWithoutOutcome = append(listOfFeaturesWithoutOutcome, f.Name)
		}
	}

	return listOfFeaturesWithoutOutcome
}

// printEpicsWithoutFeatures prints the epics without features
func getEpicsWithoutFeatures(epics []*JiraEpic, features []*JiraFeature) []string {
	listOfEpicsWithoutFeatures := []string{}
	listOfEpicsWithFeatures := []string{}
	for _, e := range epics {
		for _, f := range features {
			for _, ef := range f.Epics {
				if ef.Name == e.Name {
					listOfEpicsWithFeatures = append(listOfEpicsWithFeatures, e.Name)
				}
			}
		}
	}

	for _, e := range epics {
		// If the e.Name is not part of the listOfEpicsWithFeatures array, add it to the listOfEpicsWithoutFeatures
		if !contains(listOfEpicsWithFeatures, e.Name) {
			listOfEpicsWithoutFeatures = append(listOfEpicsWithoutFeatures, e.Name)
		}
	}

	return listOfEpicsWithoutFeatures
}

// printStoriesWithoutEpics prints the stories without epics
func getStoriesWithoutEpics(stories []*JiraStory, epics []*JiraEpic) []string {
	listOfStoriesWithoutEpics := []string{}
	listOfStoriesWithEpics := []string{}
	for _, s := range stories {
		for _, e := range epics {
			for _, es := range e.Stories {
				if es.Name == s.Name {
					listOfStoriesWithEpics = append(listOfStoriesWithEpics, s.Name)
				}
			}
		}
	}

	for _, s := range stories {
		// If the s.Name is not part of the listOfStoriesWithEpics array, add it to the listOfStoriesWithoutEpics
		if !contains(listOfStoriesWithEpics, s.Name) {
			listOfStoriesWithoutEpics = append(listOfStoriesWithoutEpics, s.Name)
		}
	}

	return listOfStoriesWithoutEpics
}

// printOutcomes without features
func getOutcomesWithoutFeatures(outcomes []*JiraOutcome) []string {
	listofOutcomesWithoutFeatures := []string{}
	// Print the Outcomes
	for _, o := range outcomes {
		if len(o.Features) == 0 {
			listofOutcomesWithoutFeatures = append(listofOutcomesWithoutFeatures, o.Name)
		}
	}

	return listofOutcomesWithoutFeatures
}

// printFeaturesWithoutEpics prints the features without epics
func getFeaturesWithoutEpics(features []*JiraFeature) []string {
	listOfFeaturesWithoutEpics := []string{}
	// Print the Features
	for _, f := range features {
		if len(f.Epics) == 0 {
			listOfFeaturesWithoutEpics = append(listOfFeaturesWithoutEpics, f.Name)
		}
	}
	return listOfFeaturesWithoutEpics
}

type JiraSubTask struct {
	Name string
}

type JiraTask struct {
	Name        string
	JiraSubTask []*JiraSubTask
}

type JiraStory struct {
	Name        string
	Description string
	Status      string
	Tasks       []*JiraTask
	Epics       []*JiraEpic
	Features    []*JiraFeature
	Outcomes    []*JiraOutcome
}

type JiraEpic struct {
	Name    string
	Stories []*JiraStory
}

type JiraFeature struct {
	Name  string
	Epics []*JiraEpic
}

type JiraOutcome struct {
	Name     string
	Features []*JiraFeature
}

func contains(arr []string, val string) bool {
	for _, a := range arr {
		if a == val {
			return true
		}
	}
	return false
}

func containsTask(arr []*JiraTask, val string) bool {
	for _, a := range arr {
		if a.Name == val {
			return true
		}
	}
	return false
}

func containsSubTask(arr []*JiraSubTask, val string) bool {
	for _, a := range arr {
		if a.Name == val {
			return true
		}
	}
	return false
}

func containsStory(arr []*JiraStory, val string) bool {
	for _, a := range arr {
		if a.Name == val {
			return true
		}
	}
	return false
}

func containsEpic(arr []*JiraEpic, val string) bool {
	for _, a := range arr {
		if a.Name == val {
			return true
		}
	}
	return false
}

func containsFeature(arr []*JiraFeature, val string) bool {
	for _, a := range arr {
		if a.Name == val {
			return true
		}
	}
	return false
}

func containsOutcome(arr []*JiraOutcome, val string) bool {
	for _, a := range arr {
		if a.Name == val {
			return true
		}
	}
	return false
}
