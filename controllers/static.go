package controllers

import (
	"html/template"
	"net/http"
)

func StaticHandler(tpl Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, nil)
	}
}

func FAQ(tpl Template) http.HandlerFunc {
	questions := []struct {
		Question string
		Answer   template.HTML
	}{
		{
			Question: "What is Pic Bucket?",
			Answer:   "Pic Bucket is a web application designed for storing and organizing photos. Users can upload, store, and manage their photos using this platform.",
		},
		{
			Question: "How do I upload photos to Pic Bucket?",
			Answer:   "To upload photos to Pic Bucket, log in to your account, go to the 'Upload' section, click on the 'Choose File' button, select the photos you want to upload from your device, and then click 'Upload'.",
		},
		{
			Question: "Can I create albums to organize my photos?",
			Answer:   "Yes, you can create albums in Pic Bucket to organize your photos. After uploading photos, you can create a new album and assign selected photos to that album.",
		},
		{
			Question: "How do I view my saved photos?",
			Answer:   "You can view your saved photos by going to the 'My Photos' or 'Albums' section in Pic Bucket. There, you will find all the photos and albums you have saved.",
		},
		{
			Question: "Is Pic Bucket free to use?",
			Answer:   "Pic Bucket offers both free and premium plans. The free plan allows limited storage space and features, while the premium plan provides more storage and additional features. You can choose the plan that best suits your needs.",
		},
		{
			Question: "How do I contact support?",
			Answer:   `Email us - <a href="mailto:support@picbucket.com">support@picbucket.com</a>`,
		},
	}
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, questions)
	}
}
