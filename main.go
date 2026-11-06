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
	// Eingaben Vorklärbecken (VB)
	LaengeVB         float64
	BreiteVB         float64
	HoeheVB          float64
	WandstaerkeVB    float64
	FundamentdickeVB float64
	AnzahlVB         float64

	// Eingaben Nachklärbecken (NB)
	DurchmesserNB    float64
	HoeheNB          float64
	WandstaerkeNB    float64
	FundamentdickeNB float64
	AnzahlNB         float64

	// Berechnete Ergebnisse
	Berechnet bool
	Fehler    string

	// Ergebnisse VB
	VolumenVB              float64
	KubaturWandVB          float64
	AussenwandschalungVB   float64
	InnenwandschalungVB    float64
	KubaturFundamentVB     float64
	FundamentschalungVB    float64
	WandbewehrungVB        float64
	BewehrungBodenplatteVB float64

	// Ergebnisse NB
	VolumenNB              float64
	KubaturWandNB          float64
	AussenwandschalungNB   float64
	InnenwandschalungNB    float64
	KubaturFundamentNB     float64
	FundamentschalungNB    float64
	WandbewehrungNB        float64
	BewehrungBodenplatteNB float64

	// Gesamtmassen
	GesamtVolumen           float64
	GesamtKubaturBeton      float64
	GesamtWandSchalungGerad float64
	GesamtWandSchalungRund  float64
	GesamtFundamentSchalung float64
	GesamtBewehrung         float64
}

const htmlTemplate = `
<!DOCTYPE html>
<html lang="de">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Massenberechnung Kläranlage</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 900px; margin: 30px auto; padding: 0 20px; line-height: 1.5; background: #f5f7fa; color: #333; }
        .card { background: white; padding: 30px; border-radius: 8px; box-shadow: 0 4px 15px rgba(0,0,0,0.05); margin-bottom: 20px; }
        h1, h2, h3 { color: #1e293b; margin-top: 0; border-bottom: 2px solid #e2e8f0; padding-bottom: 8px; }
        .section-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 25px; margin-bottom: 25px; }
        @media (max-width: 768px) { .section-grid { grid-template-columns: 1fr; } }
        .box { background: #f8fafc; padding: 20px; border-radius: 6px; border: 1px solid #e2e8f0; }
        .form-group { margin-bottom: 15px; }
        label { display: block; font-weight: bold; margin-bottom: 5px; font-size: 14px; color: #475569; }
        input[type="number"] { width: 100%; padding: 8px 12px; box-sizing: border-box; border: 1px solid #cbd5e1; border-radius: 4px; font-size: 14px; }
        button { background: #0284c7; color: white; border: none; padding: 12px 20px; cursor: pointer; font-size: 16px; border-radius: 4px; width: 100%; font-weight: bold; transition: background 0.2s; }
        button:hover { background: #0369a1; }
        .error { color: #b91c1c; font-weight: bold; background: #fee2e2; padding: 12px; border-radius: 4px; border: 1px solid #fca5a5; margin-bottom: 20px; }
        table { width: 100%; border-collapse: collapse; margin-top: 15px; margin-bottom: 25px; font-size: 14px; }
        th, td { border: 1px solid #e2e8f0; padding: 10px 12px; text-align: left; }
        th { background-color: #0f172a; color: white; }
        tr:nth-child(even) { background-color: #f8fafc; }
        .total-row { font-weight: bold; background-color: #e0f2fe !important; }
    </style>
</head>
<body>
    <div class="card">
        <h1>Massenberechnung einer Biologischen Kläranlage</h1>
        
        {{if .Fehler}}<div class="error">{{.Fehler}}</div>{{end}}

        <form method="POST">
            <div class="section-grid">
                <!-- Vorklärbecken (Rechteckig) -->
                <div class="box">
                    <h3>Vorklärbecken (Rechteck)</h3>
                    <div class="form-group">
                        <label>Innenlänge (m):</label>
                        <input type="number" step="0.01" name="laengeVB" value="{{if .Berechnet}}{{.LaengeVB}}{{else}}15.00{{end}}" required min="0">
                    </div>
                    <div class="form-group">
                        <label>Innenbreite (m):</label>
                        <input type="number" step="0.01" name="breiteVB" value="{{if .Berechnet}}{{.BreiteVB}}{{else}}6.00{{end}}" required min="0">
                    </div>
                    <div class="form-group">
                        <label>Innenhöhe (m):</label>
                        <input type="number" step="0.01" name="hoeheVB" value="{{if .Berechnet}}{{.HoeheVB}}{{else}}4.00{{end}}" required min="0">
                    </div>
                    <div class="form-group">
                        <label>Wandstärke (cm):</label>
                        <input type="number" step="0.1" name="wandstaerkeVB" value="{{if .Berechnet}}{{multiply .WandstaerkeVB 100.0}}{{else}}35.0{{end}}" required min="0">
                    </div>
                    <div class="form-group">
                        <label>Fundamentstärke (cm):</label>
                        <input type="number" step="0.1" name="fundamentdickeVB" value="{{if .Berechnet}}{{multiply .FundamentdickeVB 100.0}}{{else}}40.0{{end}}" required min="0">
                    </div>
                    <div class="form-group">
                        <label>Anzahl der Becken:</label>
                        <input type="number" name="anzahlVB" value="{{if .Berechnet}}{{.AnzahlVB}}{{else}}2{{end}}" required min="1" step="1">
                    </div>
                </div>

                <!-- Nachklärbecken (Rund) -->
                <div class="box">
                    <h3>Nachklärbecken (Rund)</h3>
                    <div class="form-group">
                        <label>Innendurchmesser (m):</label>
                        <input type="number" step="0.01" name="durchmesserNB" value="{{if .Berechnet}}{{.DurchmesserNB}}{{else}}20.00{{end}}" required min="0">
                    </div>
                    <div class="form-group">
                        <label>Innenhöhe (m):</label>
                        <input type="number" step="0.01" name="hoeheNB" value="{{if .Berechnet}}{{.HoeheNB}}{{else}}3.50{{end}}" required min="0">
                    </div>
                    <div class="form-group">
                        <label>Wandstärke (cm):</label>
                        <input type="number" step="0.1" name="wandstaerkeNB" value="{{if .Berechnet}}{{multiply .WandstaerkeNB 100.0}}{{else}}30.0{{end}}" required min="0">
                    </div>
                    <div class="form-group">
                        <label>Fundamentstärke (cm):</label>
                        <input type="number" step="0.1" name="fundamentdickeNB" value="{{if .Berechnet}}{{multiply .FundamentdickeNB 100.0}}{{else}}40.0{{end}}" required min="0">
                    </div>
                    <div class="form-group">
                        <label>Anzahl der Becken:</label>
                        <input type="number" name="anzahlNB" value="{{if .Berechnet}}{{.AnzahlNB}}{{else}}2{{end}}" required min="1" step="1">
                    </div>
                </div>
            </div>
            <button type="submit">Massen berechnen</button>
        </form>

                {{if .Berechnet}}
        <h2>Berechnungsergebnisse</h2>
        
        <h3>1. Rohbaumassen Vorklärbecken</h3>
        <table>
            <thead>
                <tr>
                    <th>Position</th>
                    <th style="text-align: right;">Wert</th>
                </tr>
            </thead>
            <tbody>
                <tr><td>Fassungsvermögen gesamt</td><td style="text-align: right;"><strong>{{printf "%.2f" .VolumenVB}} m³</strong></td></tr>
                <tr><td>Kubatur Wandbeton</td><td style="text-align: right;">{{printf "%.2f" .KubaturWandVB}} m³</td></tr>
                <tr><td>Kubatur Fundamentbeton</td><td style="text-align: right;">{{printf "%.2f" .KubaturFundamentVB}} m³</td></tr>
                <tr><td>Fläche Außenwandschalung</td><td style="text-align: right;">{{printf "%.2f" .AussenwandschalungVB}} m²</td></tr>
                <tr><td>Fläche Innenwandschalung</td><td style="text-align: right;">{{printf "%.2f" .InnenwandschalungVB}} m²</td></tr>
                <tr><td>Fläche Fundamentschalung</td><td style="text-align: right;">{{printf "%.2f" .FundamentschalungVB}} m²</td></tr>
                <tr><td>Bewehrungsstahl Wände</td><td style="text-align: right;">{{printf "%.2f" .WandbewehrungVB}} t</td></tr>
                <tr><td>Bewehrungsstahl Bodenplatten</td><td style="text-align: right;">{{printf "%.2f" .BewehrungBodenplatteVB}} t</td></tr>
            </tbody>
        </table>

        <h3>2. Rohbaumassen Nachklärbecken</h3>
        <table>
            <thead>
                <tr>
                    <th>Position</th>
                    <th style="text-align: right;">Wert</th>
                </tr>
            </thead>
            <tbody>
                <tr><td>Fassungsvermögen gesamt</td><td style="text-align: right;"><strong>{{printf "%.2f" .VolumenNB}} m³</strong></td></tr>
                <tr><td>Kubatur Wandbeton</td><td style="text-align: right;">{{printf "%.2f" .KubaturWandNB}} m³</td></tr>
                <tr><td>Kubatur Fundamentbeton</td><td style="text-align: right;">{{printf "%.2f" .KubaturFundamentNB}} m³</td></tr>
                <tr><td>Fläche Außenwandschalung (Rundschalung)</td><td style="text-align: right;">{{printf "%.2f" .AussenwandschalungNB}} m²</td></tr>
                <tr><td>Fläche Innenwandschalung (Rundschalung)</td><td style="text-align: right;">{{printf "%.2f" .InnenwandschalungNB}} m²</td></tr>
                <tr><td>Fläche Fundamentschalung</td><td style="text-align: right;">{{printf "%.2f" .FundamentschalungNB}} m²</td></tr>
                <tr><td>Bewehrungsstahl Wände</td><td style="text-align: right;">{{printf "%.2f" .WandbewehrungNB}} t</td></tr>
                <tr><td>Bewehrungsstahl Bodenplatten</td><td style="text-align: right;">{{printf "%.2f" .BewehrungBodenplatteNB}} t</td></tr>
            </tbody>
        </table>

        <h3>3. Gesamtrohbaumassen der Kläranlage</h3>
        <table>
            <thead>
                <tr>
                    <th>Gesamtposition</th>
                    <th style="text-align: right;">Masse / Einheit</th>
                </tr>
            </thead>
            <tbody>
                <tr><td>Gesamtes Fassungsvermögen</td><td style="text-align: right;"><strong>{{printf "%.2f" .GesamtVolumen}} m³</strong></td></tr>
                <tr><td>Die gesamte Kubatur des Betons</td><td style="text-align: right;"><strong>{{printf "%.2f" .GesamtKubaturBeton}} m³</strong></td></tr>
                <tr><td>Gesamte Fläche gerade Wandschalung (VB)</td><td style="text-align: right;">{{printf "%.2f" .GesamtWandSchalungGerad}} m²</td></tr>
                <tr><td>Gesamte Fläche Rundschalung (NB)</td><td style="text-align: right;">{{printf "%.2f" .GesamtWandSchalungRund}} m²</td></tr>
                <tr><td>Gesamte Fläche Fundamentschalung</td><td style="text-align: right;">{{printf "%.2f" .GesamtFundamentSchalung}} m²</td></tr>
                <tr><td>Gesamtmenge Bewehrungsstahl</td><td style="text-align: right;"><strong>{{printf "%.2f" .GesamtBewehrung}} t</strong></td></tr>

            </tbody>
        </table>
        {{end}}
    </div>
</body>
</html>
`

func main() {
	// Helper-Funktion fürs Template registrieren, um cm wieder in m umzurechnen
	funcMap := template.FuncMap{
		"multiply": func(a, b float64) float64 { return a * b },
	}
	tmpl := template.Must(template.New("klaeranlage").Funcs(funcMap).Parse(htmlTemplate))
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Diese Zeile erlaubt das Einbetten auf deiner Homepage:
		w.Header().Del("X-Frame-Options")
		
		if r.Method != http.MethodPost {
			tmpl.Execute(w, BerechnungsErgebnis{})
			return
		}

		}
		// Parsen der Web-Eingaben
		data := AnlageData{
			LaengeVB:         parseFormFloat(r.FormValue("laengeVB")),
			BreiteVB:         parseFormFloat(r.FormValue("breiteVB")),
			HoeheVB:          parseFormFloat(r.FormValue("hoeheVB")),
			WandstaerkeVB:    parseFormFloat(r.FormValue("wandstaerkeVB")) / 100.0, // cm zu m
			FundamentdickeVB: parseFormFloat(r.FormValue("fundamentdickeVB")) / 100.0,
			AnzahlVB:         parseFormFloat(r.FormValue("anzahlVB")),
			DurchmesserNB:    parseFormFloat(r.FormValue("durchmesserNB")),
			HoeheNB:          parseFormFloat(r.FormValue("hoeheNB")),
			WandstaerkeNB:    parseFormFloat(r.FormValue("wandstaerkeNB")) / 100.0, // cm zu m
			FundamentdickeNB: parseFormFloat(r.FormValue("fundamentdickeNB")) / 100.0,
			AnzahlNB:         parseFormFloat(r.FormValue("anzahlNB")),
		}
		// Konstanten aus deinem Originalcode
		const ueberstandfundament = 0.2
		const bewehrungbodenplatte = 0.12
		const bewehrungwandVB = 0.18
		const bewehrungwandNB = 0.15
		// --- MATHEMATISCHE BERECHNUNGEN (DEIN ORIGINAL-ALGORITHMUS) ---
		// Vorklärbecken
		data.VolumenVB = data.LaengeVB * data.BreiteVB * data.HoeheVB * data.AnzahlVB
		gesamtVolumenWaendeVB := (data.LaengeVB + 2*data.WandstaerkeVB) * (data.BreiteVB + 2*data.WandstaerkeVB) * data.HoeheVB
		innenVolumenWaendeVB := data.LaengeVB * data.BreiteVB * data.HoeheVB
		data.KubaturWandVB = (gesamtVolumenWaendeVB - innenVolumenWaendeVB) * data.AnzahlVB
		data.AussenwandschalungVB = (2 * ((data.LaengeVB + 2*data.WandstaerkeVB) + (data.BreiteVB + 2*data.WandstaerkeVB)) * data.HoeheVB) * data.AnzahlVB
		data.InnenwandschalungVB = (2 * (data.LaengeVB + data.BreiteVB) * data.HoeheVB) * data.AnzahlVB
		data.KubaturFundamentVB = (data.LaengeVB + 2*(data.WandstaerkeVB+ueberstandfundament)) * (data.BreiteVB + 2*(data.WandstaerkeVB+ueberstandfundament)) * data.FundamentdickeVB * data.AnzahlVB
		data.FundamentschalungVB = (2 * ((data.LaengeVB + 2*data.WandstaerkeVB) + (data.BreiteVB + 2*data.WandstaerkeVB)) * data.FundamentdickeVB) * data.AnzahlVB
		data.WandbewehrungVB = data.KubaturWandVB * bewehrungwandVB
		data.BewehrungBodenplatteVB = data.KubaturFundamentVB * bewehrungbodenplatte
		// Nachklärbecken
		rInnenNB := data.DurchmesserNB / 2.0
		rAussenNB := rInnenNB + data.WandstaerkeNB
		rFundamentNB := rAussenNB + ueberstandfundament
		data.VolumenNB = math.Pi * math.Pow(rInnenNB, 2) * data.HoeheNB * data.AnzahlNB
		data.KubaturWandNB = (math.Pi*math.Pow(rAussenNB, 2) - math.Pi*math.Pow(rInnenNB, 2)) * data.HoeheNB * data.AnzahlNB
		data.AussenwandschalungNB = (2 * math.Pi * rAussenNB * data.HoeheNB) * data.AnzahlNB
		data.InnenwandschalungNB = (2 * math.Pi * rInnenNB * data.HoeheNB) * data.AnzahlNB
		data.KubaturFundamentNB = math.Pi * math.Pow(rFundamentNB, 2) * data.FundamentdickeNB * data.AnzahlNB
		data.FundamentschalungNB = (2 * math.Pi * rFundamentNB * data.FundamentdickeNB) * data.AnzahlNB
		data.WandbewehrungNB = data.KubaturWandNB * bewehrungwandNB
		data.BewehrungBodenplatteNB = data.KubaturFundamentNB * bewehrungbodenplatte
		// Gesamtergebnisse zusammenschlagen
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
		port = "8080"
	}
	fmt.Println("Server läuft auf Port " + port)
	http.ListenAndServe(":"+port, nil)
}

// Hilfsfunktion zum sicheren Parsen der Strings aus dem Webformular
func parseFormFloat(value string) float64 {
	val, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0.0
	}
	return val
}
