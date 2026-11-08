package main

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"os"
	"strconv"
)

type AnlageData struct {
	LaengeVB               float64
	BreiteVB               float64
	HoeheVB                float64
	WandstaerkeVB          float64
	FundamentdickeVB       float64
	AnzahlVB               float64
	VolumenVB              float64
	KubaturWandVB          float64
	AussenwandschalungVB   float64
	InnenwandschalungVB    float64
	KubaturFundamentVB     float64
	FundamentschalungVB    float64
	WandbewehrungVB        float64
	BewehrungBodenplatteVB float64

	DurchmesserNB          float64
	HoeheNB                float64
	WandstaerkeNB          float64
	FundamentdickeNB       float64
	AnzahlNB               float64
	VolumenNB              float64
	KubaturWandNB          float64
	AussenwandschalungNB   float64
	InnenwandschalungNB    float64
	KubaturFundamentNB     float64
	FundamentschalungNB    float64
	WandbewehrungNB        float64
	BewehrungBodenplatteNB float64

	GesamtVolumen           float64
	GesamtKubaturBeton      float64
	GesamtWandSchalungGerad float64
	GesamtWandSchalungRund  float64
	GesamtFundamentSchalung float64
	GesamtBewehrung         float64
	Berechnet               bool
}

// Setze hier im htmlTemplate-String einfach dein komplettes, originales CSS & Layout wieder ein!
const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Massenberechnung Kläranlage</title>
    <style>
        /* HIER DEIN ORIGINALES CSS EINFÜGEN */
        .grid { display: flex; gap: 20px; }
        .box { border: 1px solid #ccc; padding: 20px; flex: 1; }
        .form-group { margin-bottom: 15px; }
        label { display: block; margin-bottom: 5px; }
    </style>
</head>
<body>
    <h2>Massenberechnung einer Kläranlage</h2>
    <form method="POST">
        <div class="grid">
            <div class="box">
                <h3>Vorklärbecken (Rechteck)</h3>
                <div class="form-group">
                    <label>Länge (m):</label>
                    <input type="number" step="0.01" name="laengeVB" placeholder="60.00" value="{{if .Berechnet}}{{printf "%.2f" .LaengeVB}}{{end}}">
                </div>
                <div class="form-group">
                    <label>Breite (m):</label>
                    <input type="number" step="0.01" name="breiteVB" placeholder="30.00" value="{{if .Berechnet}}{{printf "%.2f" .BreiteVB}}{{end}}">
                </div>
                <div class="form-group">
                    <label>Höhe (m):</label>
                    <input type="number" step="0.01" name="hoeheVB" placeholder="8.00" value="{{if .Berechnet}}{{printf "%.2f" .HoeheVB}}{{end}}">
                </div>
                <div class="form-group">
                    <label>Wandstärke (cm):</label>
                    <input type="number" step="0.1" name="wandstaerkeVB" placeholder="80.0" value="{{if .Berechnet}}{{printf "%.1f" .WandstaerkeVB}}{{end}}">
                </div>
                <div class="form-group">
                    <label>Fundamentstärke (cm):</label>
                    <input type="number" step="0.1" name="fundamentdickeVB" placeholder="90.0" value="{{if .Berechnet}}{{printf "%.1f" .FundamentdickeVB}}{{end}}">
                </div>
                <div class="form-group">
                    <label>Anzahl:</label>
                    <input type="number" name="anzahlVB" placeholder="2" value="{{if .Berechnet}}{{.AnzahlVB}}{{end}}">
                </div>
            </div>
            <div class="box">
                <h3>Nachklärbecken (Rund)</h3>
                <div class="form-group">
                    <label>Durchmesser (m):</label>
                    <input type="number" step="0.01" name="durchmesserNB" placeholder="30.00" value="{{if .Berechnet}}{{printf "%.2f" .DurchmesserNB}}{{end}}">
                </div>
                <div class="form-group">
                    <label>Höhe (m):</label>
                    <input type="number" step="0.01" name="hoeheNB" placeholder="7.00" value="{{if .Berechnet}}{{printf "%.2f" .HoeheNB}}{{end}}">
                </div>
                <div class="form-group">
                    <label>Wandstärke (cm):</label>
                    <input type="number" step="0.1" name="wandstaerkeNB" placeholder="50.0" value="{{if .Berechnet}}{{printf "%.1f" .WandstaerkeNB}}{{end}}">
                </div>
                <div class="form-group">
                    <label>Fundamentstärke (cm):</label>
                    <input type="number" step="0.1" name="fundamentdickeNB" placeholder="80.0" value="{{if .Berechnet}}{{printf "%.1f" .FundamentdickeNB}}{{end}}">
                </div>
                <div class="form-group">
                    <label>Anzahl:</label>
                    <input type="number" name="anzahlNB" placeholder="2" value="{{if .Berechnet}}{{.AnzahlNB}}{{end}}">
                </div>
            </div>
        </div>
        <button type="submit">Berechnung starten</button>
    </form>

    {{if .Berechnet}}
        <!-- HIER DEINE ORIGINALEN AUSGABEFELDER FÜR DIE ERGEBNISSE EINFÜGEN -->
    {{end}}
</body>
</html>
`

func main() {
	tmpl := template.Must(template.New("v").Parse(htmlTemplate))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Del("X-Frame-Options")

		if r.Method != http.MethodPost {
			tmpl.Execute(w, AnlageData{})
			return
		}

		// Datenstruktur speichert die reinen cm-Werte (ohne /100)
		data := AnlageData{
			LaengeVB:         p(r.FormValue("laengeVB")),
			BreiteVB:         p(r.FormValue("breiteVB")),
			HoeheVB:          p(r.FormValue("hoeheVB")),
			WandstaerkeVB:    p(r.FormValue("wandstaerkeVB")),
			FundamentdickeVB: p(r.FormValue("fundamentdickeVB")),
			AnzahlVB:         p(r.FormValue("anzahlVB")),
			DurchmesserNB:    p(r.FormValue("durchmesserNB")),
			HoeheNB:          p(r.FormValue("hoeheNB")),
			WandstaerkeNB:    p(r.FormValue("wandstaerkeNB")),
			FundamentdickeNB: p(r.FormValue("fundamentdickeNB")),
			AnzahlNB:         p(r.FormValue("anzahlNB")),
		}

		// Temporäre Meter-Variablen NUR für die mathematischen Formeln
		wStaerkeVB_m := data.WandstaerkeVB / 100.0
		fDickeVB_m := data.FundamentdickeVB / 100.0
		wStaerkeNB_m := data.WandstaerkeNB / 100.0
		fDickeNB_m := data.FundamentdickeNB / 100.0

		const ueberstand = 0.2
		const bBoden = 0.12
		const bWandVB = 0.18
		const bWandNB = 0.15

		// Berechnungen nutzen konsequent die neuen Meter-Variablen
		data.VolumenVB = data.LaengeVB * data.BreiteVB * data.HoeheVB * data.AnzahlVB
		gVolVB := (data.LaengeVB + 2*wStaerkeVB_m) * (data.BreiteVB + 2*wStaerkeVB_m) * data.HoeheVB
		iVolVB := data.LaengeVB * data.BreiteVB * data.HoeheVB
		data.KubaturWandVB = (gVolVB - iVolVB) * data.AnzahlVB
		data.AussenwandschalungVB = (2 * ((data.LaengeVB + 2*wStaerkeVB_m) + (data.BreiteVB + 2*wStaerkeVB_m)) * data.HoeheVB) * data.AnzahlVB
		data.InnenwandschalungVB = (2 * (data.LaengeVB + data.BreiteVB) * data.HoeheVB) * data.AnzahlVB
		data.KubaturFundamentVB = (data.LaengeVB + 2*(wStaerkeVB_m+ueberstand)) * (data.BreiteVB + 2*(wStaerkeVB_m+ueberstand)) * fDickeVB_m * data.AnzahlVB
		data.FundamentschalungVB = (2 * ((data.LaengeVB + 2*wStaerkeVB_m) + (data.BreiteVB + 2*wStaerkeVB_m)) * fDickeVB_m) * data.AnzahlVB
		data.WandbewehrungVB = data.KubaturWandVB * bWandVB
		data.BewehrungBodenplatteVB = data.KubaturFundamentVB * bBoden

		rIn := data.DurchmesserNB / 2.0
		rAus := rIn + wStaerkeNB_m
		rFund := rAus + ueberstand
		data.VolumenNB = math.Pi * math.Pow(rIn, 2) * data.HoeheNB * data.AnzahlNB
		data.KubaturWandNB = (math.Pi*math.Pow(rAus, 2) - math.Pi*math.Pow(rIn, 2)) * data.HoeheNB * data.AnzahlNB
		data.AussenwandschalungNB = (2 * math.Pi * rAus * data.HoeheNB) * data.AnzahlNB
		data.InnenwandschalungNB = (2 * math.Pi * rIn * data.HoeheNB) * data.AnzahlNB
		data.KubaturFundamentNB = math.Pi * math.Pow(rFund, 2) * fDickeNB_m * data.AnzahlNB
		data.FundamentschalungNB = (2 * math.Pi * rFund * fDickeNB_m) * data.AnzahlNB
		data.WandbewehrungNB = data.KubaturWandNB * bWandNB
		data.BewehrungBodenplatteNB = data.KubaturFundamentNB * bBoden

		data.GesamtVolumen = data.VolumenVB + data.VolumenNB
		data.GesamtKubaturBeton = data.KubaturWandVB + data.KubaturFundamentVB + data.KubaturWandNB + data.KubaturFundamentNB
		data.GesamtWandSchalungGerad = data.AussenwandschalungVB + data.InnenwandschalungVB
		data.GesamtWandSchalungRund = data.AussenwandschalungNB + data.InnenwandschalungNB
		data.GesamtFundamentSchalung = data.FundamentschalungVB + data.FundamentschalungNB
		data.GesamtBewehrung = data.WandbewehrungVB + data.BewehrungBodenplatteVB + data.WandbewehrungNB + data.BewehrungBodenplatteNB
		data.Berechnet = true
		tmpl.Execute(w, data)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "10000"
	}
	fmt.Println("Server läuft online...")
	http.ListenAndServe("0.0.0.0:"+port, nil)
}

func p(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}
