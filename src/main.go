package main

import ( 
	"html/template"
	"fmt"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

const mainTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Personal Blog</title>
    <link rel="stylesheet" href="style.css">
</head>
<body>
    <a href="create.html">Create New Post</a>
    <h1>Personal Blog</h1>
    <main>
{{range $i := .}}        <article>
            <a href="article{{$i}}.html"><h2>Article {{$i}}</h2></a>
        </article>
{{end}}    </main>
</body>
</html>`

const articleTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <link rel="stylesheet" href="style.css">
</head>
<body>
    <a href="index.html">Main Page</a>
    <div class='article-section'>
        <div class='article-group'>
            <h1>{{.Title}}</h1>
        </div>
        <div class='article-group'>
            <h2>{{.Date}}</h2>
        </div>
        <div class='article-group'>
            <article>
                {{.Content}}
            </article>
        </div>
    </div>
</body>
</html>`

type ArticleRequest struct {
    Title   string
    Date    string
    Content string
	Count 	int
}

type ArticleResponse struct {
	success string `json:"result"`
	Error string `json:"err,omitempty"`
}

func getArticleCount() (int, error) {
    data, err := os.ReadFile("articles.json")
    if err != nil {
        return -1, err
    }
    
    var result struct {
        TotalCount int `json:"total_count"`
    }
    
    err = json.Unmarshal(data, &result)
    if err != nil {
        return -1, err
    }
    
    return result.TotalCount, nil
}

func incrementArticleCount(title, date string) error {
    data, err := os.ReadFile("articles.json")
    if err != nil {
        return err
    }
    
    var result struct {
        TotalCount int `json:"total_count"`
        Articles   []struct {
            ID    int    `json:"id"`
            Title string `json:"title"`
            Date  string `json:"date"`
        } `json:"articles"`
    }
    
    err = json.Unmarshal(data, &result)
    if err != nil {
        return err
    }
    
    result.TotalCount++
    
    newArticle := struct {
        ID    int    `json:"id"`
        Title string `json:"title"`
        Date  string `json:"date"`
    }{
        ID:    result.TotalCount,
        Title: title,
        Date:  date,
    }
    
    result.Articles = append(result.Articles, newArticle)
    
    updatedData, err := json.Marshal(result)
    if err != nil {
        return err
    }
    
    err = os.WriteFile("articles.json", updatedData, 0644)
    return err
}

func convert(article ArticleRequest) error {
	count, countErr := getArticleCount()

	fmt.Println(count)

	if countErr != nil {
		return fmt.Errorf("Failed to get article count:", countErr)
	}

	tmpl, tempErr := template.New("article").Parse(articleTemplate)

	if tempErr != nil {
		return fmt.Errorf("failed to make template: ", tempErr)
	}

	filename := fmt.Sprintf("article%d.html", count + 1)
	file, createErr := os.Create(filename)

	if createErr != nil {
		return fmt.Errorf("Failed to create article:", createErr)
	}

	defer file.Close()

	tmpl.Execute(file, article)

	if incErr := incrementArticleCount(article.Title, article.Date)

	incErr != nil {
		return fmt.Errorf("failed to increment article:", incErr)
	}

	return nil
}

func updateMainPage() error {
    count, err := getArticleCount()
    if err != nil {
        return err
    }
    
    tmpl, err := template.New("main").Parse(mainTemplate)
    if err != nil {
        return err
    }
    
    // Create a slice of article numbers [1, 2, 3, ..., count]
    articleNumbers := make([]int, count)
    for i := 0; i < count; i++ {
        articleNumbers[i] = i + 1
    }
    
    file, err := os.Create("index.html")
    if err != nil {
        return err
    }
    defer file.Close()
    
    return tmpl.Execute(file, articleNumbers)
}

func handlePublish(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    w.Header().Set("Content-Type", "application/json") // Tell the client we're sending JSON
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ArticleRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(ArticleResponse{Error: "Invalid request"})
		return
	}

	err := convert(req)

	if err != nil {
		json.NewEncoder(w).Encode(ArticleResponse{Error: err.Error()})
		return
	}

	updateMainErr := updateMainPage()

	if updateMainErr != nil {
		json.NewEncoder(w).Encode(ArticleResponse{Error: updateMainErr.Error()})
		return
	}

	json.NewEncoder(w).Encode(ArticleResponse{success: "Success!"})
}

func main() {
	http.HandleFunc("/publish", handlePublish)

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
