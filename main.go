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
        h3 { color: #0284c7; margin-bottom: 12px; border-bottom: 1px solid #e2e8f0; padding-bottom: 4px; }
        .grid { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; margin-bottom: 20px; }
        @media (max-width: 600px) { .grid { grid-template-columns: 1fr; } }
        .box { background: #f8fafc; padding: 15px; border-radius: 6px; border: 1px solid #e2e8f0; }
        .form-group { margin-bottom: 10px; display: flex; justify-content: space-between; align-items: center; }
        label { font-size: 14px; font-weight: bold; color: #475569; }
        input[type="number"] { width: 100px; padding: 6px; border: 1px solid #cbd5e1; border-radius: 4px; text-align: right; }
        button { background: #0284c7; color: white; border: none; padding: 12px; cursor: pointer; font-size: 16px; border-radius: 4px; width: 100%; font-weight: bold; transition: background 0.2s; }
        button:hover { background: #0369a1; }
        table { width: 100%; border-collapse: collapse; margin-top: 15px; margin-bottom: 20px; font-size: 14px; }
        th, td { border: 1px solid #e2e8f0; padding: 8px 12px; text-align: left; }
        th { background-color: #0f172a; color: white; }
        tr:nth-child(even) { background-color: #f8fafc; }
        .r { text-align: right; }
    </style>
</head>
<body>
    <div class="card">
        <h2>Massenberechnung einer Biologischen Kläranlage</h2>
        <form method="POST">
            <div class="grid">
                <!-- Vorklärbecken -->
                <div class="box">
                    <h3>Vorklärbecken (Rechteck)</h3>
                    <div class="form-group"><label>Länge (m):</label><input type="number" step="0.01" name="laengeVB" value="15.00"></div>
                    <div class="form-group"><label>Breite (m):</label><input type="number" step="0.01" name="breiteVB" value="6.00"></div>
                    <div class="form-group"><label>Höhe (m):</label><input type="number" step="0.01" name="hoeheVB" value="4.00"></div>
                    <div class="form-group"><label>Wandstärke (cm):</label><input type="number" step="0.1" name="wandstaerkeVB" value="35.0"></div>
                    <div class="form-group"><label>Fundamentstärke (cm):</label><input type="number" step="0.1" name="fundamentdickeVB" value="40.0"></div>
                    <div class="form-group"><label>Anzahl:</label><input type="number" name="anzahlVB" value="2"></div>
                </div>

                <!-- Nachklärbecken -->
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
        <h2>Berechnungsergebnisse</h2>
        
        <h3>1. Rohbaumassen Vorklärbecken</h3>
        <table>
            <tr><td>Fassungsvermögen gesamt</td><td class="r"><strong>{{printf "%.2f" .VolumenVB}} m³</strong></td></tr>
            <tr><td>Kubatur Wandbeton</td><td class="r">{{printf "%.2f" .KubaturWandVB}} m³</td></tr>
            <tr><td>Kubatur Fundamentbeton</td><td class="r">{{printf "%.2f" .KubaturFundamentVB}} m³</td></tr>
            <tr><td>Fläche Außenwandschalung</td><td class="r">{{printf "%.2f" .AussenwandschalungVB}} m²</td></tr>
            <tr><td>Fläche Innenwandschalung</td><td class="r">{{printf "%.2f" .InnenwandschalungVB}} m²</td></tr>
            <tr><td>Fläche Fundamentschalung</td><td class="r">{{printf "%.2f" .FundamentschalungVB}} m²</td></tr>
            <tr><td>Bewehrungsstahl Wände</td><td class="r">{{printf "%.2f" .WandbewehrungVB}} t</td></tr>
            <tr><td>Bewehrungsstahl Bodenplatten</td><td class="r">{{printf "%.2f" .BewehrungBodenplatteVB}} t</td></tr>
        </table>

        <h3>2. Rohbaumassen Nachklärbecken</h3>
        <table>
            <tr><td>Fassungsvermögen gesamt</td><td class="r"><strong>{{printf "%.2f" .VolumenNB}} m³</strong></td></tr>
            <tr><td>Kubatur Wandbeton</td><td class="r">{{printf "%.2f" .KubaturWandNB}} m³</td></tr>
            <tr><td>Kubatur Fundamentbeton</td><td class="r">{{printf "%.2f" .KubaturFundamentNB}} m³</td></tr>
            <tr><td>Fläche Außenwandschalung (Rund)</td><td class="r">{{printf "%.2f" .AussenwandschalungNB}} m²</td></tr>
            <tr><td>Fläche Innenwandschalung (Rund)</td><td class="r">{{printf "%.2f" .InnenwandschalungNB}} m²</td></tr>
            <tr><td>Fläche Fundamentschalung</td><td class="r">{{printf "%.2f" .FundamentschalungNB}} m²</td></tr>
            <tr><td>Bewehrungsstahl Wände</td><td class="r">{{printf "%.2f" .WandbewehrungNB}} t</td></tr>
            <tr><td>Bewehrungsstahl Bodenplatten</td><td class="r">{{printf "%.2f" .BewehrungBodenplatteNB}} t</td></tr>
        </table>

        <h3>3. Gesamtrohbaumassen der Kläranlage</h3>
        <table>
            <tr><th>Gesamtposition</th><th class="r">Masse / Einheit</th></tr>
            <tr><td>Gesamtes Fassungsvermögen</td><td class="r"><strong>{{printf "%.2f" .GesamtVolumen}} m³</strong></td></tr>
            <tr><td>Die gesamte Kubatur des Betons</td><td class="r"><strong>{{printf "%.2f" .GesamtKubaturBeton}} m³</strong></td></tr>
            <tr><td>Gesamte Fläche gerade Wandschalung (VB)</td><td class="r">{{printf "%.2f" .GesamtWandSchalungGerad}} m²</td></tr>
            <tr><td>Gesamte Fläche Rundschalung (NB)</td><td class="r">{{printf "%.2f" .GesamtWandSchalungRund}} m²</td></tr>
            <tr><td>Gesamte Fläche Fundamentschalung</td><td class="r">{{printf "%.2f" .GesamtFundamentSchalung}} m²</td></tr>
            <tr><td>Gesamtmenge Bewehrungsstahl</td><td class="r"><strong>{{printf "%.2f" .GesamtBewehrung}} t</strong></td></tr>
        </table>
        {{end}}
    </div>
</body>
</html>
`
