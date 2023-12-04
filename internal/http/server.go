package http

import (
	"net/http"

	"github.com/Nebezdar/tech_lvl0.git/internal/app"
)

type Server struct {
	orderService *app.OrderService
}

func NewServer(orderService *app.OrderService) *Server {
	return &Server{orderService: orderService}
}

func (s *Server) Start(port string) {
	http.HandleFunc("/order", s.orderHandler)
	http.ListenAndServe(":"+port, nil)
}

func (s *Server) orderHandler(w http.ResponseWriter, r *http.Request) {
	orderID := r.URL.Query().Get("id")
	if orderID == "" {
		http.Error(w, "Order ID is required", http.StatusBadRequest)
		return
	}

	data, err := s.orderService.GetOrderByID(orderID)
	if err != nil {
		http.Error(w, "Error retrieving order data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
