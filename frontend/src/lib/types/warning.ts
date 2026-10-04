/** Mirrors the backend's WarningResponse (app/models/warning.go): a DHMZ warning via Meteoalarm. */
export interface Warning {
  level: "yellow" | "orange" | "red";
  /** wind, snow-ice, thunderstorm, fog, high-temperature, low-temperature, rain, … */
  type: string;
  /** Croatian headline, e.g. "Žuto upozorenje za vjetar". */
  event: string;
  description: string;
  onset: string;
  expires: string;
}
