package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type PageData struct {
	Berechnet                  bool
	Fehler                     string
	LaengeVB                   float64
	BreiteVB                   float64
	HoeheVB                    float64
	WandstaerkeVB              float64
	FundamentdickeVB           float64
	AnzahlVB                   int
	DurchmesserNB              float64
	HoeheNB                    float64
	WandstaerkeNB              float64
	FundamentdickeNB           float64
	AnzahlNB                   int
	VolumenVB                  float64
	KubaturWandVB              float64
	KubaturFundamentVB         float64
	AussenwandschalungVB       float64
	InnenwandschalungVB        float64
	FundamentschalungVB        float64
	WandbewehrungVB            float64
	BewehrungBodenplatteVB     float64
	VolumenNB                  float64
	KubaturWandNB              float64
	KubaturFundamentNB         float64
	AussenwandschalungNB       float64
	InnenwandschalungNB        float64
	FundamentschalungNB        float64
	WandbewehrungNB            float64
	BewehrungBodenplatteNB     float64
	GesamtVolumen              float64
	GesamtKubaturBeton         float64
	GesamtWandSchalungGerad    float64
	GesamtWandSchalungRund     float64
	GesamtFundamentSchalung    float64
	GesamtBewehrung            float64
}

func main() {
	htmlInhalt := strings.Join([]string{
		`<!DOCTYPE html>`,
		`<html lang="de">`,
		`<head>`,
		`    <meta charset="UTF-8">`,
		`    <meta name="viewport" content="width=device-width, initial-scale=1.0">`,
		`    <title>Berechnungsprogramm Kläranlage</title>`,
		`    <style>`,
		`        body { font-family: 'Roboto', sans-serif; max-width: 800px; margin: 20px auto; padding: 20px; background: #000000; color: #333; }`,
		`        .card { background: white; padding: 25px; border-radius: 5px; box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05); }`,
		`        h1 { color: #111; border-bottom: 2px solid #eee; padding-bottom: 8px; margin-top: 0; font-size: 24px; }`,
		`        h2 { color: #111; margin-top: 20px; font-size: 20px; border-bottom: 1px solid #e2e8f0; padding-bottom: 4px; }`,
		`        h3 { color: #428787; margin-bottom: 12px; border-bottom: 1px solid #e2e8f0; padding-bottom: 4px; font-size: 18px; }`,
		`        .grid { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; margin-bottom: 20px; }`,
		`        @media (max-width: 600px) { .grid { grid-template-columns: 1fr; } }`,
		`        .box { background: #f8fafc; padding: 15px; border-radius: 5px; border: 1px solid #e2e8f0; }`,
		`        .form-group { margin-bottom: 10px; display: flex; justify-content: space-between; align-items: center; }`,
		`        label { font-size: 14px; font-weight: bold; color: #475569; }`,
		`        input[type="number"] { width: 100px; padding: 6px; border: 1px solid #cbd5e1; border-radius: 4px; text-align: right; font-size: 16px; }`,
		`        button { background: #428787; color: white; border: none; padding: 12px; cursor: pointer; font-size: 16px; border-radius: 4px; width: 100%; font-weight: bold; transition: background 0.2s; }`,
		`        button:hover:not(:disabled) { background: #336666; }`,
		`        button:disabled { background: #cbd5e1 !important; color: #94a3b8 !important; cursor: not-allowed !important; opacity: 0.7; }`,
		`        table { table-layout: fixed; width: 100%; border-collapse: collapse; margin-top: 15px; margin-bottom: 20px; font-size: 14px; }`,
		`        th, td { border: 1px solid #e2e8f0; padding: 8px 12px; text-align: left; }`,
		`        th { background-color: #428787 !important; color: white !important; }`,
		`        .table-vb tr:nth-child(even), .table-nb tr:nth-child(even) { background-color: #f8fafc; }`,
		`        .table-gesamt tr:nth-child(odd) { background-color: #f8fafc; }`,
		`        .r { text-align: right; }`,
		`        table th:first-child, table td:first-child { width: 60% !important; box-sizing: border-box; }`,
		`        table th:last-child, table td:last-child { width: 40% !important; box-sizing: border-box; }`,
		`        /* Zwingt das Hellblau auf die gesamte Tabellenzeile inklusive Zellen */`,
		`        table.table-gesamt tr.final-highlight td { background-color: #e2f2f2 !important; font-weight: bold !important; }`,
		`    </style>`,
		`    <script>`,
		`        function pruefeFormular() {`,
		`            var form = document.getElementById("calcForm");`,
		`            var submitBtn = document.getElementById("submitBtn");`,
		`            if (form && submitBtn) {`,
		`                if (form.checkValidity()) { submitBtn.removeAttribute("disabled"); }`,
		`                else { submitBtn.setAttribute("disabled", "true"); }`,
		`            }`,
		`        }`,
		`        document.addEventListener("DOMContentLoaded", function() {`,
		`            var form = document.getElementById("calcForm");`,
		`            if (form) { form.addEventListener("input", pruefeFormular); }`,
		`            pruefeFormular();`,
		`        });`,
		`    </script>`,
		`</head>`,
		`<body>`,
		`    <div class="card">`,
		`        <h1>Massenberechnung Kläranlage</h1>`,
		`        <form method="POST" id="calcForm">`,
		`            <div class="grid">`,
		`                <div class="box">`,
		`                    <h3>Vorklärbecken (Rechteck)</h3>`,
		`                    <div class="form-group"><label>Länge (m):</label><input type="number" step="0.01" name="laengeVB" placeholder="00.00" value="{{if .Berechnet}}{{printf \"%.2f\" .LaengeVB}}{{end}}" required></div>`,
		`                    <div class="form-group"><label>Breite (m):</label><input type="number" step="0.01" name="breiteVB" placeholder="00.00" value="{{if .Berechnet}}{{printf \"%.2f\" .BreiteVB}}{{end}}" required></div>`,
		`                    <div class="form-group"><label>Höhe (m):</label><input type="number" step="0.01" name="hoeheVB" placeholder="00.00" value="{{if .Berechnet}}{{printf \"%.2f\" .HoeheVB}}{{end}}" required></div>`,
		`                    <div class="form-group"><label>Wandstärke (cm):</label><input type="number" step="0.1" name="wandstaerkeVB" placeholder="00.00" value="{{if .Berechnet}}{{printf \"%.1f\" .WandstaerkeVB}}{{end}}" required></div>`,
		`                    <div class="form-group"><label>Fundamentstärke (cm):</label><input type="number" step="0.1" name="fundamentdickeVB" placeholder="00.00" value="{{if .Berechnet}}{{printf \"%.1f\" .FundamentdickeVB}}{{end}}" required></div>`,
		`                    <div class="form-group"><label>Anzahl:</label><input type="number" name="anzahlVB" placeholder="0" value="{{if .Berechnet}}{{.AnzahlVB}}{{end}}" required></div>`,
		`                </div>`,
		`                <div class="box">`,
		`                    <h3>Nachklärbecken (Rund)</h3>`,
		`                    <div class="form-group"><label>Durchmesser (m):</label><input type="number" step="0.01" name="durchmesserNB" placeholder="00.00" value="{{if .Berechnet}}{{printf \"%.2f\" .DurchmesserNB}}{{end}}" required></div>`,
		`                    <div class="form-group"><label>Höhe (m):</label><input type="number" step="0.01" name="hoeheNB" placeholder="00.00" value="{{if .Berechnet}}{{printf \"%.2f\" .HoeheNB}}{{end}}" required></div>`,
		`                    <div class="form-group"><label>Wandstärke (cm):</label><input type="number" step="0.1" name="wandstaerkeNB" placeholder="00.00" value="{{if .Berechnet}}{{printf \"%.1f\" .WandstaerkeNB}}{{end}}" required></div>`,
		`                    <div class="form-group"><label>Fundamentstärke (cm):</label><input type="number" step="0.1" name="fundamentdickeNB" placeholder="00.00" value="{{if .Berechnet}}{{printf \"%.1f\" .FundamentdickeNB}}{{end}}" required></div>`,
		`                    <div class="form-group"><label>Anzahl:</label><input type="number" name="anzahlNB" placeholder="0" value="{{if .Berechnet}}{{.AnzahlNB}}{{end}}" required></div>`,
		`                </div>`,
		`            </div>`,
		`            <div><h6>&copy; Ing. Burkhard Kaar - Studienprojekt imperative Programmierung, Berechnung der Rohbaumassen einer Kläranlage</h6></div>`,
		`            <button type="submit" id="submitBtn" disabled>Berechnung starten</button>`,
		`        </form>`,
		`        {{if .Berechnet}}`,
		`        <h2 class="result-title">Berechnungsergebnisse</h2>`,
		`        <div class="grid">`,
		`            <div class="box">`,
		`                <h3>Rohbaumassen Vorklärbecken</h3>`,
		`                <table class="table-vb">`,
		`                    <thead><tr><th>Position</th><th class="r">Wert</th></tr></thead>`,
		`                    <tbody>`,
		`                        <tr><td>Fassungsvermögen gesamt</td><td class="r">{{printf \"%.2f\" .VolumenVB}} m³</td></tr>`,
		`                        <tr><td>Kubatur Wandbeton</td><td class="r">{{printf \"%.2f\" .KubaturWandVB}} m³</td></tr>`,
		`                        <tr><td>Kubatur Fundamentbeton</td><td class="r">{{printf \"%.2f\" .KubaturFundamentVB}} m³</td></tr>`,
		`                        <tr><td>Fläche Außenwandschalung</td><td class="r">{{printf \"%.2f\" .AussenwandschalungVB}} m²</td></tr>`,
		`                        <tr><td>Fläche Innenwandschalung</td><td class="r">{{printf \"%.2f\" .InnenwandschalungVB}} m²</td></tr>`,
		`                        <tr><td>Fläche Fundamentschalung</td><td class="r">{{printf \"%.2f\" .FundamentschalungVB}} m²</td></tr>`,
		`                        <tr><td>Bewehrungsstahl Wände</td><td class="r">{{printf \"%.2f\" .WandbewehrungVB}} t</td></tr>`,
		`                        <tr><td>Bewehrungsstahl Bodenplatten</td><td class="r">{{printf \"%.2f\" .BewehrungBodenplatteVB}} t</td></tr>`,
		`                    </tbody>`,
		`                </table>`,
		`            </div>`,
		`            <div class="box">`,
		`                <h3>Rohbaumassen Nachklärbecken</h3>`,
		`                <table class="table-nb">`,
		`                    <thead><tr><th>Position</th><th class="r">Wert</th></tr></thead>`,
		`                    <tbody>`,
		`                        <tr><td>Fassungsvermögen gesamt</td><td class="r">{{printf "%.2f" .VolumenNB}} m³</td></tr>`,
		`                        <tr><td>Kubatur Wandbeton</td><td class="r">{{printf "%.2f" .KubaturWandNB}} m³</td></tr>`,
		`                        <tr><td>Kubatur Fundamentbeton</td><td class="r">{{printf "%.2f" .KubaturFundamentNB}} m³</td></tr>`,
		`                        <tr><td>Fläche Außenwandschalung (Rund)</td><td class="r">{{printf "%.2f" .AussenwandschalungNB}} m²</td></tr>`,
		`                        <tr><td>Fläche Innenwandschalung (Rund)</td><td class="r">{{printf "%.2f" .InnenwandschalungNB}} m²</td></tr>`,
		`                        <tr><td>Fläche Fundamentschalung</td><td class="r">{{printf "%.2f" .FundamentschalungNB}} m²</td></tr>`,
		`                        <tr><td>Bewehrungsstahl Wände</td><td class="r">{{printf "%.2f" .WandbewehrungNB}} t</td></tr>`,
		`                        <tr><td>Bewehrungsstahl Bodenplatten</td><td class="r">{{printf "%.2f" .BewehrungBodenplatteNB}} t</td></tr>`,
		`                    </tbody>`,
		`                </table>`,
		`            </div>`,
		`        </div>`,
		`        <h3>Gesamtmassen Projekt</h3>`,
		`        <table class="table-gesamt">`,
		`            <thead><tr><th>Position</th><th class="r">Gesamtwert</th></tr></thead>`,
		`            <tbody>`,
		`                <tr><td>Gesamtes Fassungsvermögen</td><td class="r">{{printf "%.2f" .GesamtVolumen}} m³</td></tr>`,
		`                <tr><td>Gesamte Kubatur Beton</td><td class="r">{{printf "%.2f" .GesamtKubaturBeton}} m³</td></tr>`,
		`                <tr><td>Gesamte Wandschalung (Gerade)</td><td class="r">{{printf "%.2f" .GesamtWandSchalungGerad}} m²</td></tr>`,
		`                <tr><td>Gesamte Wandschalung (Rund)</td><td class="r">{{printf "%.2f" .GesamtWandSchalungRund}} m²</td></tr>`,
		`                <tr><td>Gesamte Fundamentschalung</td><td class="r">{{printf "%.2f" .GesamtFundamentSchalung}} m²</td></tr>`,
		`                <tr class="final-highlight"><td>Gesamter Bewehrungsstahl</td><td class="r">{{printf "%.2f" .GesamtBewehrung}} t</td></tr>`,
		`            </tbody>`,
		`        </table>`,
		`        {{end}}`,
		`    </div>`,
		`</body>`,
		`</html>`,
}, "\n")
tmpl := template.Must(template.New("index").Parse(htmlInhalt))
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
data := PageData{}
if r.Method == http.MethodPost {
// Hilfsfunktion zur String-Konvertierung
pf := func(name string) float64 {
v, _ := strconv.ParseFloat(strings.ReplaceAll(r.FormValue(name), ",", "."), 64)
return v
}
pi := func(name string) int {
v, _ := strconv.Atoi(r.FormValue(name))
return v
}
data.LaengeVB = pf("laengeVB")
data.BreiteVB = pf("breiteVB")
data.HoeheVB = pf("hoeheVB")
data.WandstaerkeVB = pf("wandstaerkeVB")
data.FundamentdickeVB = pf("fundamentdickeVB")
data.AnzahlVB = pi("anzahlVB")
data.DurchmesserNB = pf("durchmesserNB")
data.HoeheNB = pf("hoeheNB")
data.WandstaerkeNB = pf("wandstaerkeNB")
data.FundamentdickeNB = pf("fundamentdickeNB")
data.AnzahlNB = pi("anzahlNB")
// Dummy-Kalkulationen (Hier setzt deine Logik an)
data.VolumenVB = data.LaengeVB * data.BreiteVB * data.HoeheVB * float64(data.AnzahlVB)
data.KubaturWandVB = 45.5
data.KubaturFundamentVB = 32.1
data.VolumenNB = 3.1415 * (data.DurchmesserNB / 2) * (data.DurchmesserNB / 2) * data.HoeheNB * float64(data.AnzahlNB)
data.KubaturWandNB = 28.4
data.KubaturFundamentNB = 19.8
data.GesamtVolumen = data.VolumenVB + data.VolumenNB
data.GesamtKubaturBeton = data.KubaturWandVB + data.KubaturFundamentVB + data.KubaturWandNB + data.KubaturFundamentNB
data.GesamtBewehrung = 14.85
data.Berechnet = true
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
