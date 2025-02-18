package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const imageDir = "./images" // โฟลเดอร์ที่ใช้เก็บภาพ

// ตรวจสอบว่าโฟลเดอร์ที่เก็บภาพมีอยู่หรือไม่ ถ้าไม่มีก็สร้างขึ้น
func init() {
	if _, err := os.Stat(imageDir); os.IsNotExist(err) {
		err := os.MkdirAll(imageDir, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating image directory:", err)
			os.Exit(1)
		}
	}
}

// ฟังก์ชันอัพโหลดภาพ
func uploadImage(w http.ResponseWriter, r *http.Request) {
	// ตรวจสอบว่าเป็นคำขอแบบ POST หรือไม่
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// อ่านไฟล์ที่อัพโหลด
	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Failed to read image", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// สร้างชื่อไฟล์ใหม่
	fileName := fmt.Sprintf("%s.jpg", strings.ReplaceAll(fmt.Sprintf("%d", r.ContentLength), ".", ""))

	// สร้างไฟล์ใหม่ในโฟลเดอร์ที่เก็บภาพ
	outFile, err := os.Create(filepath.Join(imageDir, fileName))
	if err != nil {
		http.Error(w, "Failed to save image", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	// คัดลอกข้อมูลจากไฟล์ที่อัพโหลดไปยังไฟล์ในเซิร์ฟเวอร์
	_, err = io.Copy(outFile, file)
	if err != nil {
		http.Error(w, "Failed to copy image", http.StatusInternalServerError)
		return
	}

	// ส่ง URL ของภาพที่อัพโหลดกลับไป
	imageURL := fmt.Sprintf("http://localhost:8080/images/%s", fileName)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Image uploaded successfully. Access it at: %s", imageURL)))
}

// ฟังก์ชันดาวน์โหลดภาพ
func getImage(w http.ResponseWriter, r *http.Request) {
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
