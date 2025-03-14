package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const imageDir = "./images"

func init() {
	if _, err := os.Stat(imageDir); os.IsNotExist(err) {
		err := os.MkdirAll(imageDir, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating image directory:", err)
			os.Exit(1)
		}
	}
}

// Add CORS headers
func setCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")                              // Allow all domains (for testing purposes)
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")            // Allow certain methods
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Custom-Header") // Allow specific headers
}

func uploadImage(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	setCORSHeaders(w)

	// Handle preflight requests for CORS (OPTIONS method)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// ตรวจสอบว่าเป็นคำขอแบบ POST หรือไม่
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// อ่านไฟล์ที่อัปโหลด
	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Failed to read image", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// ใช้ชื่อไฟล์จากที่อัปโหลด
	fileName := handler.Filename

	// สร้างไฟล์ใหม่ในโฟลเดอร์ที่เก็บภาพ
	outFile, err := os.Create(filepath.Join(imageDir, fileName))
	if err != nil {
		http.Error(w, "Failed to save image", http.StatusInternalServerError)
		return
	}

	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		http.Error(w, "Failed to copy image", http.StatusInternalServerError)
		return
	}

	// ส่ง JSON Response ที่มีเฉพาะชื่อไฟล์
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"filename": fileName})
}

func getImage(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	setCORSHeaders(w)

	// รับชื่อไฟล์จาก URL
	imageName := strings.TrimPrefix(r.URL.Path, "/images/")

	// ตรวจสอบว่าไฟล์มีอยู่ในโฟลเดอร์หรือไม่
	imagePath := filepath.Join(imageDir, imageName)
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		http.Error(w, "Image not found", http.StatusNotFound)
		return
	}

	// อ่านไฟล์และส่งกลับ
	http.ServeFile(w, r, imagePath)
}

func main() {
	http.HandleFunc("/upload", uploadImage) // API สำหรับอัพโหลดภาพ
	http.HandleFunc("/images/", getImage)   // API สำหรับดาวน์โหลดภาพ

	fmt.Println("Image storage server is running at http://localhost:8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
