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
	LaengeVB, BreiteVB, HoeheVB, WandstaerkeVB, FundamentdickeVB, AnzahlVB float64
	DurchmesserNB, HoeheNB, WandstaerkeNB, FundamentdickeNB, AnzahlNB      float64
	Berechnet                                                              bool
	VolumenVB, KubaturWandVB, AussenwandschalungVB, InnenwandschalungVB    float64
	KubaturFundamentVB, FundamentschalungVB, WandbewehrungVB, BewehrungBodenplatteVB float64
	VolumenNB, KubaturWandNB, AussenwandschalungNB, InnenwandschalungNB    float64
	KubaturFundamentNB, FundamentschalungNB, WandbewehrungNB, BewehrungBodenplatteNB float64
	GesamtVolumen, GesamtKubaturBeton, GesamtWandSchalungGerad             float64
	GesamtWandSchalungRund, GesamtFundamentSchalung, GesamtBewehrung       float64
}

const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Massenberechnung Kläranlage</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 800px; margin: 20px auto; padding: 20px; background: #fafafa; color: #333; }
        .card { background: white; padding: 25px; border-radius: 8px; box-shadow: 0 4px 12px rgba(0,0,0,0.05); }
        h2 { color: #111; border-bottom: 2px solid #eee; padding-bottom: 8px; margin-top: 0; }
        .result-title { margin-top: 40px !important; }
        h3 { color: #428787; margin-bottom: 12px; border-bottom: 1px solid #e2e8f0; padding-bottom: 4px; }
        .grid { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; margin-bottom: 20px; }
        @media (max-width: 600px) { .grid { grid-template-columns: 1fr; } }
        .box { background: #f8fafc; padding: 15px; border-radius: 6px; border: 1px solid #e2e8f0; }
        .form-group { margin-bottom: 10px; display: flex; justify-content: space-between; align-items: center; }
        label { font-size: 14px; font-weight: bold; color: #475569; }
        input[type="number"] { width: 100px; padding: 6px; border: 1px solid #cbd5e1; border-radius: 4px; text-align: right; }
        button { background: #428787; color: white; border: none; padding: 12px; cursor: pointer; font-size: 16px; border-radius: 4px; width: 100%; font-weight: bold; transition: background 0.2s; }
        button:hover { background: #336666; }
        table { width: 100%; border-collapse: collapse; margin-top: 15px; margin-bottom: 20px; font-size: 14px; }
        th, td { border: 1px solid #e2e8f0; padding: 8px 12px; text-align: left; }
        th { background-color: #428787 !important; color: white !important; }
        .table-vb tr:nth-child(even), .table-nb tr:nth-child(even) { background-color: #f8fafc; }
        .table-gesamt tr:nth-child(odd) { background-color: #f8fafc; }
        .r { text-align: right; }
    </style>
</head>
<body>
    <div class="card">
        <h2>Massenberechnung einer Biologischen Kläranlage</h2>
        <form method="POST">
            <div class="grid">
                <div class="box">
                    <h3>Vorklärbecken (Rechteck)</h3>
                    <div class="form-group"><label>Länge (m):</label><input type="number" step="0.01" name="laengeVB" value="15.00"></div>
                    <div class="form-group"><label>Breite (m):</label><input type="number" step="0.01" name="breiteVB" value="6.00"></div>
                    <div class="form-group"><label>Höhe (m):</label><input type="number" step="0.01" name="hoeheVB" value="4.00"></div>
                    <div class="form-group"><label>Wandstärke (cm):</label><input type="number" step="0.1" name="wandstaerkeVB" value="35.0"></div>
                    <div class="form-group"><label>Fundamentstärke (cm):</label><input type="number" step="0.1" name="fundamentdickeVB" value="40.0"></div>
                    <div class="form-group"><label>Anzahl:</label><input type="number" name="anzahlVB" value="2"></div>
                </div>
                <div class="box">
                    <h3>Nachklärbecken (Rund)</h3>
                    <div class="form-group"><label>Durchmesser (m):</label><input type="number" step="0.01" name="durchmesserNB" value="20.00"></div>
                    <div class="form-group"><label>Höhe (m):</label><input type="number" step="0.01" name="hoeheNB" value="3.50"></div>
                    <div class="form-group"><label>Wandstärke (cm):</label><input type="number" step="0.1" name="wandstaerkeNB" value="30.0"></div>
                    <div class="form-group"><label>Fundamentstärke (cm):</label><input type="number" step="0.1" name="fundamentdickeNB" value="40.0"></div>
                    <div class="form-group"><label>Anzahl:</label><input type="number" name="anzahlNB" value="2"></div>
                </div>
            </div>
            <button type="submit">Massen berechnen</button>
        </form>

        {{if .Berechnet}}
        <h2 class="result-title">Berechnungsergebnisse</h2>
        
        <h3>1. Rohbaumassen Vorklärbecken</h3>
        <table class="table-vb">
            <thead>
                <tr>
                    <th>Position</th>
                    <th class="r">Wert</th>
                </tr>
            </thead>
            <tbody>
                <tr><td>Fassungsvermögen gesamt</td><td class="r">{{printf "%.2f" .VolumenVB}} m³</td></tr>
                <tr><td>Kubatur Wandbeton</td><td class="r">{{printf "%.2f" .KubaturWandVB}} m³</td></tr>
                <tr><td>Kubatur Fundamentbeton</td><td class="r">{{printf "%.2f" .KubaturFundamentVB}} m³</td></tr>
                <tr><td>Fläche Außenwandschalung</td><td class="r">{{printf "%.2f" .AussenwandschalungVB}} m²</td></tr>
                <tr><td>Fläche Innenwandschalung</td><td class="r">{{printf "%.2f" .InnenwandschalungVB}} m²</td></tr>
                <tr><td>Fläche Fundamentschalung</td><td class="r">{{printf "%.2f" .FundamentschalungVB}} m²</td></tr>
                <tr><td>Bewehrungsstahl Wände</td><td class="r">{{printf "%.2f" .WandbewehrungVB}} t</td></tr>
                <tr><td>Bewehrungsstahl Bodenplatten</td><td class="r">{{printf "%.2f" .BewehrungBodenplatteVB}} t</td></tr>
            </tbody>
        </table>

        <h3>2. Rohbaumassen Nachklärbecken</h3>
        <table class="table-nb">
            <thead>
                <tr>
                    <th>Position</th>
                    <th class="r">Wert</th>
                </tr>
            </thead>
            <tbody>
                <tr><td>Fassungsvermögen gesamt</td><td class="r">{{printf "%.2f" .VolumenNB}} m³</td></tr>
                <tr><td>Kubatur Wandbeton</td><td class="r">{{printf "%.2f" .KubaturWandNB}} m³</td></tr>
                <tr><td>Kubatur Fundamentbeton</td><td class="r">{{printf "%.2f" .KubaturFundamentNB}} m³</td></tr>
                <tr><td>Fläche Außenwandschalung (Rund)</td><td class="r">{{printf "%.2f" .AussenwandschalungNB}} m²</td></tr>
                <tr><td>Fläche Innenwandschalung (Rund)</td><td class="r">{{printf "%.2f" .InnenwandschalungNB}} m²</td></tr>
                <tr><td>Fläche Fundamentschalung</td><td class="r">{{printf "%.2f" .FundamentschalungNB}} m²</td></tr>
                <tr><td>Bewehrungsstahl Wände</td><td class="r">{{printf "%.2f" .WandbewehrungNB}} t</td></tr>
                <tr><td>Bewehrungsstahl Bodenplatten</td><td class="r">{{printf "%.2f" .BewehrungBodenplatteNB}} t</td></tr>
            </tbody>
        </table>

        <h3>3. Gesamtrohbaumassen der Kläranlage</h3>
        <table class="table-gesamt">
            <thead>
                <tr>
                    <th>Gesamtposition</th>
                    <th class="r">Masse / Einheit</th>
                </tr>
            </thead>
            <tbody>
                <tr><td>Gesamtes Fassungsvermögen</td><td class="r">{{printf "%.2f" .GesamtVolumen}} m³</td></tr>
                <tr><td>Die gesamte Kubatur des Betons</td><td class="r">{{printf "%.2f" .GesamtKubaturBeton}} m³</td></tr>
                <tr><td>Gesamte Fläche gerade Wandschalung (VB)</td><td class="r">{{printf "%.2f" .GesamtWandSchalungGerad}} m²</td></tr>
                <tr><td>Gesamte Fläche Rundschalung (NB)</td><td class="r">{{printf "%.2f" .GesamtWandSchalungRund}} m²</td></tr>
                <tr><td>Gesamte Fläche Fundamentschalung</td><td class="r">{{printf "%.2f" .GesamtFundamentSchalung}} m²</td></tr>
                <tr><td>Gesamtmenge Bewehrungsstahl</td><td class="r">{{printf "%.2f" .GesamtBewehrung}} t</td></tr>
            </tbody>
        </table>
        {{end}}
    </div>
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

		data := AnlageData{
			LaengeVB:         p(r.FormValue("laengeVB")),
			BreiteVB:         p(r.FormValue("breiteVB")),
			HoeheVB:          p(r.FormValue("hoeheVB")),
			WandstaerkeVB:    p(r.FormValue("wandstaerkeVB")) / 100.0,
			FundamentdickeVB: p(r.FormValue("fundamentdickeVB")) / 100.0,
			AnzahlVB:         p(r.FormValue("anzahlVB")),
			DurchmesserNB:    p(r.FormValue("durchmesserNB")),
			HoeheNB:          p(r.FormValue("hoeheNB")),
			WandstaerkeNB:    p(r.FormValue("wandstaerkeNB")) / 100.0,
			FundamentdickeNB: p(r.FormValue("fundamentdickeNB")) / 100.0,
			AnzahlNB:         p(r.FormValue("anzahlNB")),
		}

		const ueberstand = 0.2
		const bBoden = 0.12
		const bWandVB = 0.18
		const bWandNB = 0.15

		data.VolumenVB = data.LaengeVB * data.BreiteVB * data.HoeheVB * data.AnzahlVB
		gVolVB := (data.LaengeVB + 2*data.WandstaerkeVB) * (data.BreiteVB + 2*data.WandstaerkeVB) * data.HoeheVB
		iVolVB := data.LaengeVB * data.BreiteVB * data.HoeheVB
		data.KubaturWandVB = (gVolVB - iVolVB) * data.AnzahlVB
		data.AussenwandschalungVB = (2 * ((data.LaengeVB + 2*data.WandstaerkeVB) + (data.BreiteVB + 2*data.WandstaerkeVB)) * data.HoeheVB) * data.AnzahlVB
		data.InnenwandschalungVB = (2 * (data.LaengeVB + data.BreiteVB) * data.HoeheVB) * data.AnzahlVB
		data.KubaturFundamentVB = (data.LaengeVB + 2*(data.WandstaerkeVB+ueberstand)) * (data.BreiteVB + 2*(data.WandstaerkeVB+ueberstand)) * data.FundamentdickeVB * data.AnzahlVB
		data.FundamentschalungVB = (2 * ((data.LaengeVB + 2*data.WandstaerkeVB) + (data.BreiteVB + 2*data.WandstaerkeVB)) * data.FundamentdickeVB) * data.AnzahlVB
		data.WandbewehrungVB = data.KubaturWandVB * bWandVB
		data.BewehrungBodenplatteVB = data.KubaturFundamentVB * bBoden

		rIn := data.DurchmesserNB / 2.0
		rAus := rIn + data.WandstaerkeNB
		rFund := rAus + ueberstand
		data.VolumenNB = math.Pi * math.Pow(rIn, 2) * data.HoeheNB * data.AnzahlNB
		data.KubaturWandNB = (math.Pi*math.Pow(rAus, 2) - math.Pi*math.Pow(rIn, 2)) * data.HoeheNB * data.AnzahlNB
		data.AussenwandschalungNB = (2 * math.Pi * rAus * data.HoeheNB) * data.AnzahlNB
		data.InnenwandschalungNB = (2 * math.Pi * rIn * data.HoeheNB) * data.AnzahlNB
		data.KubaturFundamentNB = math.Pi * math.Pow(rFund, 2) * data.FundamentdickeNB * data.AnzahlNB
		data.FundamentschalungNB = (2 * math.Pi * rFund * data.FundamentdickeNB) * data.AnzahlNB
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

