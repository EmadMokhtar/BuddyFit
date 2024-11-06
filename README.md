# BuddyFit

BuddyFit is an AI-powered, conversational interface designed to help users build and maintain a consistent exercise routine.

## Architecture

![architecture](./arch.png)

## Components

- Fetcher: Fetches YouTube videos' transcript and convert them into a full text format.
- Importer: Imports the full text into the TimeScaleDB database.
- BuddyFit CLI: Command line interface for BuddyFit, that answers user's questions about the exercises.

## AI Features

BuddyFit is relaying on the AI features provided by the TimescaleDB database.
- The `pgai` extension is used to run the AI queries. 
- The `pgai Vectorizer` is used to automatically build embeddings for the text data existing in TimescaleDB.

BuddyFit is using OpenAI's `text-embedding-3-small` and `gpt-4o-mini` models to generate the embeddings 
for the user's questions and the exercises' descriptions.

## Installation


```bash
go install github.com/EmadMokhtar/BuddyFit/cmd/buddyfit@v1.0.0
```

Make sure the cli is installed by running the following command:

```bash
$ buddyfit --help
Usage of buddyfit:
  -p string
    	Alias for prompt
  -prompt string
    	Prompt for the AI
````

## Usage

### Requirements

- Go 1.22+ [docs](https://golang.org/doc/install)
- [Environment Variables](#environment-variables)

#### Environment Variables

- `OPENAI_API_KEY`: The OpenAI API key [docs](https://help.openai.com/en/articles/4936850-where-do-i-find-my-openai-api-key).
- `BF_DB_URL`: The TimescaleDB database URL [docs](https://docs.timescale.com/use-timescale/latest/integrations/query-admin/about-connecting/). 

### BuddyFit CLI

After [installing](#installation) the CLI, you can run the CLI app but invoking the `buddyfit` command and pass the `-p` or `-prompt` flag with your question.

```bash
$ buddyfit -p "I want to maximize my muscle gain and going to the gym 4 times per week. Can you build a plan for me? I don't want to spend long time in the gym"
````

## Example Usage

#### Example 1

```bash
$ buddyfit -p "I want to maximize my muscle gain and going to the gym 4 times per week. Can you build a plan for me? I don't want to spend long time in the gym"
```

  Sure! Here's a sample 4-day upper and lower body split that focuses on
  maximizing
  muscle gain while keeping your gym sessions efficient. Each workout should take
  about 45-60 minutes.

  ### Day 1: Upper Body (Focus on Chest & Shoulders)

1. Bench Press: 3 sets of 8-10 reps
2. Incline Dumbbell Press: 3 sets of 8-10 reps
3. Overhead Dumbbell Press: 3 sets of 8-10 reps
4. Dumbbell Lateral Raises: 3 sets of 12-15 reps
5. Tricep Dips: 3 sets of 8-10 reps

  ### Day 2: Lower Body (Focus on Quads & Hamstrings)

1. Squats: 3 sets of 8-10 reps
2. Romanian Deadlifts: 3 sets of 8-10 reps
3. Leg Press: 3 sets of 8-10 reps
4. Lunges: 3 sets of 10 reps per leg
5. Calf Raises: 3 sets of 12-15 reps

  ### Day 3: Upper Body (Focus on Back & Arms)

1. Pull-Ups or Lat Pulldowns: 3 sets of 8-10 reps
2. Bent Over Rows: 3 sets of 8-10 reps
3. Face Pulls: 3 sets of 12-15 reps
4. Barbell Bicep Curls: 3 sets of 10-12 reps
5. Skull Crushers: 3 sets of 10-12 reps

  ### Day 4: Lower Body (Focus on Glutes & Hamstrings)

1. Deadlifts: 3 sets of 8-10 reps
2. Glute Bridges: 3 sets of 10-12 reps
3. Leg Curls: 3 sets of 10-12 reps
4. Step-Ups: 3 sets of 10 reps per leg
5. Seated Calf Raises: 3 sets of 12-15 reps

### Notes:

• Always warm up before starting your workouts.
• Aim to increase weights gradually while maintaining proper form.
• Incorporate rest intervals of 60-90 seconds between sets.
• Consider deloading every 4-6 weeks as needed to prevent fatigue.
• Adjust exercises based on your comfort and experience level.

This plan should help you build muscle efficiently without lengthy gym sessions.


#### Example 2

```bash
$ buddyfit -p "I want to maximize my muscle gain and going to the gym 4 times per week. Can you build a plan for me? I don't want to spend long time in the gym"
```


Sure! Here's a sample 4-day upper and lower body split that focuses on
maximizing muscle gain while keeping your gym sessions efficient. Each workout should take
about 45-60 minutes.

### Day 1: Upper Body (Focus on Chest & Shoulders)

1. Bench Press: 3 sets of 8-10 reps
2. Incline Dumbbell Press: 3 sets of 8-10 reps
3. Overhead Dumbbell Press: 3 sets of 8-10 reps
4. Dumbbell Lateral Raises: 3 sets of 12-15 reps
5. Tricep Dips: 3 sets of 8-10 reps

### Day 2: Lower Body (Focus on Quads & Hamstrings)

1. Squats: 3 sets of 8-10 reps
2. Romanian Deadlifts: 3 sets of 8-10 reps
3. Leg Press: 3 sets of 8-10 reps
4. Lunges: 3 sets of 10 reps per leg
5. Calf Raises: 3 sets of 12-15 reps

### Day 3: Upper Body (Focus on Back & Arms)

1. Pull-Ups or Lat Pulldowns: 3 sets of 8-10 reps
2. Bent Over Rows: 3 sets of 8-10 reps
3. Face Pulls: 3 sets of 12-15 reps
4. Barbell Bicep Curls: 3 sets of 10-12 reps
5. Skull Crushers: 3 sets of 10-12 reps

### Day 4: Lower Body (Focus on Glutes & Hamstrings)

1. Deadlifts: 3 sets of 8-10 reps
2. Glute Bridges: 3 sets of 10-12 reps
3. Leg Curls: 3 sets of 10-12 reps
4. Step-Ups: 3 sets of 10 reps per leg
5. Seated Calf Raises: 3 sets of 12-15 reps

### Notes:

- Always warm up before starting your workouts. 
- Aim to increase weights gradually while maintaining proper form. 
- Incorporate rest intervals of 60-90 seconds between sets. 
- Consider deloading every 4-6 weeks as needed to prevent fatigue. 
- Adjust exercises based on your comfort and experience level.

This plan should help you build muscle efficiently without lengthy gym sessions.


#### Example 3

```bash
$ buddyfit -p "can you make a plan to target a bigger biceps?"
```

 To target bigger biceps, you can follow this workout plan that emphasizes 
 exercises to develop both heads of the biceps and enhance their peak. Here's a 
 suggested weekly plan:

  Biceps Workout Plan:

  Day 1: Focus on Elbow Flexion and Supination

1. Barbell Curl - 4 sets of 8-12 reps
   - Focus on keeping your elbows close to your body.
2. Hammer Curl - 3 sets of 10-12 reps
   - This targets the brachialis and adds size to your biceps.
3. Concentration Curl - 3 sets of 10-12 reps
   - Ensure proper form to maximize peak contraction at the top.


  Day 2: Supination Emphasis

1. Incline Dumbbell Curl - 4 sets of 8-10 reps
   - Lying back helps emphasize the length of the biceps.
2. Reverse Curl - 3 sets of 10-12 reps
   - This targets the forearms and biceps, promoting overall arm growth.
3. Zottman Curl - 3 sets of 10-12 reps
   - Performs both supination and pronation to hit different muscle fibers.


  Day 3: Complete Arm Day (including biceps)

1. Close-Grip Bench Press - 4 sets of 8-10 reps
   - While primarily a triceps exercise, it also engages the biceps.
2. Cable Curl - 3 sets of 10-15 reps
   - Provides constant tension throughout the movement.
3. Fat Gripz or Grip Training - 3 sets of 10-12 reps
   - Improves grip strength, which helps in all lifting.


Day 4: Recovery & Stretching

- Focus on active recovery through light cardio and stretching, especially for
the arms.

Tips for Progress:

- Aim to progressively overload the muscles by increasing weight or reps over time.
- Ensure you're using proper form to effectively target the biceps and avoid injury.
- Don't neglect proper nutrition, including enough protein to support muscle growth.
- Rest and recovery are crucial; allow at least 48 hours before targeting the  same muscle group again.

By maintaining this plan and focusing on targeted exercises, you should see an
increase in bicep size and peak over time.