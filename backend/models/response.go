package models

type CommandType struct {
	Name string `json:"name"`
	Args struct {
		URL      string `json:"url"`
		Question string `json:"question"`
		Input    string `json:"input"`
		Reason   string `json:"reason"`
	} `json:"args"`
}

type MemorizeType struct {
	Subject     string `json:"subject"`
	Information string `json:"information"`
}
type ThoughtsType struct {
	Text      string   `json:"text"`
	Reasoning string   `json:"reasoning"`
	Plan      []string `json:"plan"`
	Criticism string   `json:"criticism"`
	Speak     string   `json:"speak"`
	Keywords  string   `json:"keywords"`
}

type Response struct {
	Thoughts ThoughtsType `json:"thoughts"`
	Command  CommandType  `json:"command"`
	Memorize MemorizeType `json:"memorize"`
}

/*
{
    "thoughts": {
        "text": "I have found the website Goodreads which has a list of the best hard science fiction books, and reviews from users. I should browse the website and extract the required information to create a list of top 5 recommendations to give to the user.",
        "reasoning": "Goodreads has a great list of best hard science fiction books with reviews from users, browsing their website means I can get comprehensive and reliable information.",
        "plan": "- Browse Goodreads website to find best selling hard science fiction books with good reviews\n- Collect relevant information such as book title, author and reviews\n- Create a list of five books with a combination of good reviews and rankings",
        "criticism": "I should be mindful not to generate a list based solely on top 5 rankings, but instead provide a mix of different author's work and ensure that I do not favor a particular style or sub-genre of hard science fiction books.",
        "speak": "I have found a list of the best hard science fiction books on Goodreads. I will now browse their website to gather more information for your top 5 recommendations."
    },
    "command": {
        "name": "browse_website",
        "args": {
            "url": "https://www.goodreads.com/list/show/50.The_Best_Hard_Science_Fiction_Books",
            "question": "best hard science fiction books with good reviews"
        }
    }
}
*/
