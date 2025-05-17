package api

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"time"
)

type Server struct {
	addr      string
	startTime time.Time
	gemini    *GeminiClient
	qaService *QAService
}

type healthPayload struct {
	Status string `json:"status"`
	Uptime string `json:"uptime"`
}

type QuestionRequest struct {
	Question string `json:"question"`
	Mode     string `json:"mode,omitempty"`
}

type QuestionResponse struct {
	Answer string `json:"answer"`
}

func NewServer(addr string, qaSvc *QAService) (*Server, error) {
	client, err := NewGeminiClient()
	if err != nil {
		return nil, err
	}
	return &Server{
		addr:      addr,
		startTime: time.Now(),
		gemini:    client,
		qaService: qaSvc,
	}, nil
}

func (s *Server) Run() error {
	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "hello, %q", html.EscapeString(r.URL.Path))
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf("Error Writing Response: %v ", err)
		}
	})

	router.HandleFunc("/ui", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html") // Make sure file exists
	})

	router.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		uptime := time.Since(s.startTime).String()
		healthyPayloadStatus := healthPayload{Status: "ok", Uptime: uptime}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(healthyPayloadStatus); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf("Error Writing Response: %v ", err)
		}
	})

	router.HandleFunc("/question", s.handleQuestion)
	router.HandleFunc("/history", s.handleHistory)

	log.Printf("Server has started %s", s.addr)
	return http.ListenAndServe(s.addr, router)
}

func (s *Server) handleQuestion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req QuestionRequest
	// 2. Decode Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON Payload", http.StatusBadRequest)
		return
	}
	// default to simple answer mode
	if req.Mode == "" {
		req.Mode = "answer"
	}

	// FOR NOW: stub userID=1 until i add auth middleware
	userID := 1

	answer, err := s.qaService.AskAndSave(r.Context(), userID, req.Question, req.Mode)
	if err != nil {
		log.Printf("Error in AskAndSave: %v", err)
		http.Error(w, "Internal Server Error", http.StatusServiceUnavailable)
		return
	}

	// 4. Build and write JSON response
	resp := QuestionResponse{Answer: answer}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Error writing JSON", http.StatusInternalServerError)
	}
}

func (s *Server) handleHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	userID := 1 // temporary static userID
	messages, err := s.qaService.msgRepo.ListMessages(userID)
	if err != nil {
		log.Printf("Error Loading Messages: %v", err)
		http.Error(w, "Failed to load messages", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		http.Error(w, "Encoding Failed", http.StatusInternalServerError)
	}
}
