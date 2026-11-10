package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type StromRechnung struct {
	VerbrauchKWh        float64
	EnergieArbeitspreis float64
	EnergieKostenNetto  float64
	EnergieGrundpreis   float64
	NetzArbeitspreis    float64
	NetzVerlustNetto    float64
	NetzGrundpreis      float64
	ElektroAbgabeNetto  float64
	NettoGesamt         float64
	Umsatzsteuer        float64
	BruttoGesamt        float64
	MonatsAbschlag      float64
	EffektiverKWhPreis  float64
}

type PageData struct {
	Berechnet   bool
	Fehler      string
	VerbrauchIn string
	PreisIn     string
	Ergebnis    StromRechnung
}

func berechneStrompreis(verbrauch float64, preisProKWh float64) StromRechnung {
	energieGrundpreis := 59.90
	netzArbeitspreisProKWh := 9.11
	netzVerlustProKWh := 0.40
	netzGrundpreis := 210.00
	elektroAbgabeProKWh := 1.50

	energieKostenNetto := verbrauch * (preisProKWh / 100.0)
	netzArbeitNetto := verbrauch * (netzArbeitspreisProKWh / 100.0)
	netzVerlustNetto := verbrauch * (netzVerlustProKWh / 100.0)
	abgabeNetto := verbrauch * (elektroAbgabeProKWh / 100.0)

	nettoGesamt := energieKostenNetto + energieGrundpreis + netzArbeitNetto + netzVerlustNetto + netzGrundpreis + abgabeNetto
	ust := nettoGesamt * 0.20
	bruttoGesamt := nettoGesamt + ust
	monatsAbschlag := bruttoGesamt / 12.0
	effektiverPreis := (bruttoGesamt / verbrauch) * 100.0

	return StromRechnung{
		VerbrauchKWh:        verbrauch,
		EnergieArbeitspreis: preisProKWh,
		EnergieKostenNetto:  energieKostenNetto,
		EnergieGrundpreis:   energieGrundpreis,
		NetzArbeitspreis:    netzArbeitNetto,
		NetzVerlustNetto:    netzVerlustNetto,
		NetzGrundpreis:      netzGrundpreis,
		ElektroAbgabeNetto:  abgabeNetto,
		NettoGesamt:         nettoGesamt,
		Umsatzsteuer:        ust,
		BruttoGesamt:        bruttoGesamt,
		MonatsAbschlag:      monatsAbschlag,
		EffektiverKWhPreis:  effektiverPreis,
	}
}

func main() {
	// Sicheres Zusammenfügen des HTML-Templates ohne Code-Zerstückelung
	htmlInhalt := strings.Join([]string{
		`<!DOCTYPE html>`,
		`<html lang="de">`,
		`<head>`,
		`    <meta charset="UTF-8">`,
		`    <meta name="viewport" content="width=device-width, initial-scale=1.0">`,
		`    <title>Strompreisrechner Österreich</title>`,
		`    <style>`,
		`        body { font-family: 'Roboto', sans-serif; max-width: 800px; margin: 20px auto; padding: 20px; background: #000000; color: #333; }`,
		`        .card { background: white; padding: 25px; border-radius: 5px; box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05); }`,
		`        h1 { color: #111; border-bottom: 2px solid #eee; padding-bottom: 8px; margin-top: 0; font-size: 24px; }`,
		`        h2 { color: #111; margin-top: 20px; font-size: 20px; border-bottom: 1px solid #e2e8f0; padding-bottom: 4px; }`,
		`        h3 { color: #428787; margin-bottom: 12px; border-bottom: 1px solid #e2e8f0; padding-bottom: 4px; font-size: 18px; }`,
		`        .grid { display: grid; grid-template-columns: 1fr; gap: 20px; margin-bottom: 20px; }`,
		`        .box { background: #f8fafc; padding: 15px; border-radius: 5px; border: 1px solid #e2e8f0; }`,
		`        .form-group { margin-bottom: 15px; }`,
		`        label { display: block; font-size: 14px; font-weight: bold; color: #475569; margin-bottom: 6px; }`,
		`        input[type="text"] { width: 100%; padding: 8px; border: 1px solid #cbd5e1; font-size: 16px; box-sizing: border-box; }`,
		`        button { background: #428787; color: white; border: none; padding: 12px; cursor: pointer; font-size: 16px; border-radius: 4px; width: 100%; font-weight: bold; transition: all 0.2s; }`,
		`        button:hover:not(:disabled) { background: #336666; }`,
		`        button:disabled { background: #cbd5e1; color: #94a3b8; cursor: not-allowed; opacity: 0.7; }`,
		`        .error { color: #b91c1c; font-weight: bold; background: #fee2e2; padding: 12px; border-radius: 4px; border: 1px solid #fca5a5; margin-bottom: 20px; }`,
		`        table { width: 100%; border-collapse: collapse; margin-top: 15px; margin-bottom: 20px; font-size: 14px; table-layout: fixed; }`,
		`        th, td { border: 1px solid #e2e8f0; padding: 8px 12px; text-align: left; }`,
		`        th { background-color: #428787 !important; color: white !important; }`,
		`        tr:nth-child(even) { background-color: #f8fafc; }`,
		`        .r { text-align: right; }`,
		`        .hint { color: #64748b; font-size: 12px; margin-top: 4px; display: block; font-weight: normal; }`,
		`        th:first-child, td:first-child { width: 65%; }`,
		`        th:last-child, td:last-child { white-space: nowrap; }`,
		`        .final-row td { background-color: #e2f2f2 !important; font-weight: bold !important; }`,
		`    </style>`,
		`</head>`,
		`<body>`,
		`    <div class="card">`,
		`        <h1>Strompreisrechner Österreich</h1>`,
		`        <p style="color: #428787; font-style: italic; margin-top: -5px; margin-bottom: 20px; font-size: 14px;">Berechnungsbasis: Stadtwerke Köflach / Netz Steiermark</p>`,
		`        {{if .Fehler}}<div class="error">{{.Fehler}}</div>{{end}}`,
		`        <form method="POST" id="stromForm">`,
		`            <div class="grid">`,
		`                <div class="box">`,
		`                    <h3>Eingabeparameter</h3>`,
		`                    <div class="form-group">`,
		`                        <label for="verbrauch">Jahresverbrauch (in kWh):</label>`,
		`                        <input type="text" id="verbrauch" name="verbrauch" placeholder="z.B. 3945" value="{{.VerbrauchIn}}" required>`,
		`                    </div>`,
		`                    <div class="form-group">`,
		`                        <label for="preis">Energie-Arbeitspreis (in Cent/kWh):</label>`,
		`                        <input type="text" id="preis" name="preis" placeholder="z.B. 15.90" value="{{.PreisIn}}" required>`,
		`                        <span class="hint">Reiner Energiepreis laut Tarifblatt exkl. Netz und Steuern</span>`,
		`                    </div>`,
		`                    <div><h6>&copy; Ing. Burkhard Kaar - Studienprojekt imperative Programmierung, Strompreisrechner für Österreich</h6></div>`,
		`                    <button type="submit" id="submitBtn" disabled>Berechnung starten</button>`,
		`                    {{if .Berechnet}}`,
		`                    <h2>Berechnungsergebnisse</h2>`,
		`                    <table>`,
		`                        <thead>`,
		`                            <tr><th>Kostenbestandteil</th><th class="r">Betrag (Netto)</th></tr>`,
		`                        </thead>`,
		`                        <tbody>`,
		`                            <tr><td>Reine Energiekosten (variabel)</td><td class="r">{{printf "%.2f" .Ergebnis.EnergieKostenNetto}} €</td></tr>`,
		`                            <tr><td>Energie Grundpreis (fix)</td><td class="r">{{printf "%.2f" .Ergebnis.EnergieGrundpreis}} €</td></tr>`,
		`                            <tr><td>Netznutzung (variabel)</td><td class="r">{{printf "%.2f" .Ergebnis.NetzArbeitspreis}} €</td></tr>`,
		`                            <tr><td>Netzverlustentgelt (variabel)</td><td class="r">{{printf "%.2f" .Ergebnis.NetzVerlustNetto}} €</td></tr>`,
		`                            <tr><td>Netz-Grundpreis &amp; Messpauschale (fix)</td><td class="r">{{printf "%.2f" .Ergebnis.NetzGrundpreis}} €</td></tr>`,
		`                            <tr><td>Elektrizitätsabgabe</td><td class="r">{{printf "%.2f" .Ergebnis.ElektroAbgabeNetto}} €</td></tr>`,
		`                            <tr><td>Zwischensumme Netto</td><td class="r">{{printf "%.2f" .Ergebnis.NettoGesamt}} €</td></tr>`,
		`                            <tr><td>+ 20% Umsatzsteuer</td><td class="r">{{printf "%.2f" .Ergebnis.Umsatzsteuer}} €</td></tr>`,
		`                            <tr style="background-color: #428787; color: white;"><td style="color: white;"><strong>JAHRESBETRAG BRUTTO</strong></td><td class="r" style="color: white;"><strong>{{printf "%.2f" .Ergebnis.BruttoGesamt}} €</strong></td></tr>`,
		`                            <tr><td>Monatlicher Teilzahlungsbetrag</td><td class="r">{{printf "%.2f" .Ergebnis.MonatsAbschlag}} €</td></tr>`,
		`                            <tr class="final-row"><td>Effektiver Gesamtpreis je kWh (Brutto)</td><td class="r">{{printf "%.2f" .Ergebnis.EffektiverKWhPreis}} Cent/kWh</td></tr>`,
		`                        </tbody>`,
		`                    </table>`,
		`                    {{end}}`,
		`                </div>`,
		`            </div>`,
		`        </form>`,
		`    </div>`,
		`    <script>`,
		`        function pruefeFormular() {`,
		`            var form = document.getElementById("stromForm");`,
		`            var submitBtn = document.getElementById("submitBtn");`,
		`            if (form && submitBtn) {`,
		`                if (form.checkValidity()) { submitBtn.removeAttribute("disabled"); }`,
		`                else { submitBtn.setAttribute("disabled", "true"); }`,
		`            }`,
		`        }`,
		`        var f = document.getElementById("stromForm");`,
		`        if (f) { f.addEventListener("input", pruefeFormular); }`,
		`        pruefeFormular();`,
		`    </script>`,
		`</body>`,
		`</html>`,
	}, "\n")

	tmpl := template.Must(template.New("index").Parse(htmlInhalt))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := PageData{VerbrauchIn: "", PreisIn: ""}
		if r.Method == http.MethodPost {
			vStr := strings.TrimSpace(strings.ReplaceAll(r.FormValue("verbrauch"), ",", "."))
			pStr := strings.TrimSpace(strings.ReplaceAll(r.FormValue("preis"), ",", "."))
			data.VerbrauchIn = vStr
			data.PreisIn = pStr
			verbrauch, errV := strconv.ParseFloat(vStr, 64)
			preis, errP := strconv.ParseFloat(pStr, 64)
			if errV != nil || errP != nil || verbrauch <= 0 || preis <= 0 {
				data.Fehler = "Bitte gib gültige positive Zahlen ein."
			} else {
				data.Ergebnis = berechneStrompreis(verbrauch, preis)
				data.Berechnet = true
			}
		}
		tmpl.Execute(w, data)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server läuft auf Port %s ...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
