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
<html lang="de">
<head>
    <meta charset="UTF-8">
    <title>Cross-Verfahren Berechnung</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 800px; margin: 20px auto; padding: 20px; background: #fafafa; color: #333; }
        .card { background: white; padding: 25px; border-radius: 8px; box-shadow: 0 4px 12px rgba(0,0,0,0.05); }
        h1 { color: #111; border-bottom: 2px solid #eee; padding-bottom: 8px; margin-top: 0; font-size: 24px; }
        h2 { color: #111; margin-top: 30px; font-size: 20px; border-bottom: 1px solid #e2e8f0; padding-bottom: 4px; }
        h3 { color: #428787; margin-bottom: 12px; border-bottom: 1px solid #e2e8f0; padding-bottom: 4px; font-size: 18px; }
        .grid { display: grid; grid-template-columns: 1fr; gap: 20px; margin-bottom: 20px; }
        .box { background: #f8fafc; padding: 15px; border-radius: 6px; border: 1px solid #e2e8f0; }
        .form-group { margin-bottom: 15px; }
        label { display: block; font-size: 14px; font-weight: bold; color: #475569; margin-bottom: 6px; }
        input[type="number"], input[type="text"] { width: 100%; padding: 8px; border: 1px solid #cbd5e1; border-radius: 4px; font-size: 14px; box-sizing: border-box; }
        button { background: #428787; color: white; border: none; padding: 12px; cursor: pointer; font-size: 16px; border-radius: 4px; width: 100%; font-weight: bold; transition: background 0.2s; }
        button:hover { background: #336666; }
        .error { color: #b91c1c; font-weight: bold; background: #fee2e2; padding: 12px; border-radius: 4px; border: 1px solid #fca5a5; margin-bottom: 20px; }
        table { width: 100%; border-collapse: collapse; margin-top: 15px; margin-bottom: 20px; font-size: 14px; }
        th, td { border: 1px solid #e2e8f0; padding: 8px 12px; text-align: left; }
        th { background-color: #428787; color: white; } /* Tabellenkopf in das passende Webspace-Blau geändert */
        tr:nth-child(even) { background-color: #f8fafc; }
        .r { text-align: right; }
        .hint { color: #64748b; font-size: 12px; margin-top: 4px; display: block; font-weight: normal; }
    </style>
</head>
<body>
    <div class="card">
        <h1>Durchlaufträger-Berechnung (Cross-Verfahren)</h1>
        <p style="color: #64748b; font-style: italic; margin-top: -5px; margin-bottom: 20px; font-size: 14px;">Systembelastung: Konstante Gleichlast q = 10.0 kN/m</p>

        {{if .Fehler}}
            <div class="error">{{.Fehler}}</div>
        {{end}}

        <form method="POST">
            <div class="grid">
                <div class="box">
                    <h3>Eingabeparameter Träger</h3>
                    <div class="form-group">
                        <label for="felder">Anzahl der Felder:</label>
                        <input type="number" id="felder" name="felder" value="{{if .AnzahlFelder}}{{.AnzahlFelder}}{{else}}3{{end}}" min="1" required>
                    </div>
                    <div class="form-group">
                        <label for="laengen">Feldlängen in Metern (mit Komma getrennt):</label>
                        <input type="text" id="laengen" name="laengen" value="{{.LaengenInput}}" placeholder="z.B. 4.5, 6.0, 4.5" required>
                        <span class="hint">Gib die Längen der Felder der Reihe nach ein (z.B. 5, 4.5, 5)</span>
                    </div>
                    <button type="submit">Berechnung starten</button>

                    {{if .Berechnet}}
                        <h2>Berechnungsergebnisse</h2>
                        <table>
                            <thead>
                                <tr>
                                    <th>Stabseite</th>
                                    <th class="r">µ</th>
                                    <th class="r">Moment (kNm)</th>
                                </tr>
                            </thead>
                            <tbody>
                                {{range .Seiten}}
                                <tr>
                                    <td><strong>{{.Name}}</strong></td>
                                    <td class="r">{{printf "%.2f" .Mu}}</td>
                                    <td class="r">{{printf "%.2f" .MTotal}}</td>
                                </tr>
                                {{end}}
                            </tbody>
                        </table>
                    {{end}}
                </div>
            </div>
        </form>
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

