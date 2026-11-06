package main

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"os"
	"strconv"
)

// Struktur für die Formulardaten und Berechnungen
type AnlageData struct {
	LaengeVB, BreiteVB, HoeheVB, WandstaerkeVB, FundamentdickeVB, AnzahlVB float64
	DurchmesserNB, HoeheNB, WandstaerkeNB, FundamentdickeNB, AnzahlNB      float64
	Berechnet                                                              bool
	Fehler                                                                 string
	VolumenVB, KubaturWandVB, AussenwandschalungVB, InnenwandschalungVB    float64
	KubaturFundamentVB, FundamentschalungVB, WandbewehrungVB, BewehrungBodenplateVB float64
	VolumenNB, KubaturWandNB, AussenwandschalungNB, InnenwandschalungNB    float64
	KubaturFundamentNB, FundamentschalungNB, WandbewehrungNB, BewehrungBodenplatteNB float64
	GesamtVolumen, GesamtKubaturBeton, GesamtWandSchalungGerad             float64
	GesamtWandSchalungRund, GesamtFundamentSchalung, GesamtBewehrung       float64
}

const htmlTemplate = `...` // Vollständiger HTML/CSS-Template-Code für das Webinterface der Massenberechnung

func main() {
	funcMap := template.FuncMap{
		"multiply": func(a, b float64) float64 { return a * b },
	}

	tmpl := template.Must(template.New("klaeranlage").Funcs(funcMap).Parse(htmlTemplate))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Del("X-Frame-Options")

		if r.Method != http.MethodPost {
			tmpl.Execute(w, AnlageData{})
			return
		}

		data := AnlageData{
			LaengeVB:         parseFormFloat(r.FormValue("laengeVB")),
			BreiteVB:         parseFormFloat(r.FormValue("breiteVB")),
			HoeheVB:          parseFormFloat(r.FormValue("hoeheVB")),
			WandstaerkeVB:    parseFormFloat(r.FormValue("wandstaerkeVB")) / 100.0,
			FundamentdickeVB: parseFormFloat(r.FormValue("fundamentdickeVB")) / 100.0,
			AnzahlVB:         parseFormFloat(r.FormValue("anzahlVB")),
			DurchmesserNB:    parseFormFloat(r.FormValue("durchmesserNB")),
			HoeheNB:          parseFormFloat(r.FormValue("hoeheNB")),
			WandstaerkeNB:    parseFormFloat(r.FormValue("wandstaerkeNB")) / 100.0,
			FundamentdickeNB: parseFormFloat(r.FormValue("fundamentdickeNB")) / 100.0,
			AnzahlNB:         parseFormFloat(r.FormValue("anzahlNB")),
		}

		const ueberstandfundament = 0.2
		const bewehrungbodenplatte = 0.12
		const bewehrungwandVB = 0.18
		const bewehrungwandNB = 0.15

		// Berechnungen Vorklärbecken & Nachklärbecken sowie Gesamtmassen
		data.VolumenVB = data.LaengeVB * data.BreiteVB * data.HoeheVB * data.AnzahlVB
		// ... (Berechnungslogik für Kubaturen, Schalungen und Bewehrungen)
		data.Berechnet = true
		tmpl.Execute(w, data)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Server läuft auf Port " + port)
	http.ListenAndServe(":"+port, nil)
}

func parseFormFloat(value string) float64 {
	val, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0.0
	}
	return val
}
