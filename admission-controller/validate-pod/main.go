package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func handleValidate(w http.ResponseWriter, r *http.Request) {
	var admissionReview admissionv1.AdmissionReview
	// Default: allow request
	allowed := true
	resultMsg := "Allowed: Created via controller"
		
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &admissionReview); err != nil {
		http.Error(w, "could not decode admission review", http.StatusBadRequest)
		return
	}
	if admissionReview.Request.Kind.Kind == "Pod" {
		fmt.Println("Received a Pod admission request")
		var pod map[string]interface{}
		if err := json.Unmarshal(admissionReview.Request.Object.Raw, &pod); err == nil {}
			metadata := pod["metadata"].(map[string]interface{})
			if ownerRefs, ok := metadata["ownerReferences"]; !ok || len(ownerRefs.([]interface{})) == 0 {
				allowed = false
				resultMsg = "Denied: Creating naked pods directly is not allowed"
			} }
	// Send response back to Kubernetes
	admissionReview.Response = &admissionv1.AdmissionResponse{
		UID:     admissionReview.Request.UID,
		Allowed: allowed,
		Result: &metav1.Status{
			Message: resultMsg,
		},
	}
	resp, err := json.Marshal(admissionReview)
	if err != nil {
		http.Error(w, "could not encode response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func main() {
	http.HandleFunc("/validate", handleValidate)
	fmt.Println(" Starting Pod Validator Admission Controller on :8443")
	log.Fatal(http.ListenAndServeTLS(":8443", "/tls/tls.crt", "/tls/tls.key", nil))
}
