/** Mirrors the backend's ForecastResponse (app/models/forecast.go). Times are ISO 8601 instants,
 *  all carrying the offset in force when Open-Meteo was asked (even across a DST change), so show
 *  them with Forecast.timezone, never by reading the offset; dates are YYYY-MM-DD. */
export interface CurrentForecast {
  time: string;
  /** °C */
  temperature: number;
  /** °C */
  apparentTemperature: number;
  /** % */
  humidity: number;
  /** °C */
  dewPoint: number;
  /** mm */
  precipitation: number;
  /** WMO weather interpretation code */
  weatherCode: number;
  /** % */
  cloudCover: number;
  /** hPa at sea level */
  pressure: number;
  /** km/h */
  windSpeed: number;
  /** Degrees, where the wind blows from */
  windDirection: number;
  /** km/h */
  windGusts: number;
  /** m */
  visibility: number;
  uvIndex: number;
  isDay: boolean;
}

export interface HourlyForecast {
  time: string;
  temperature: number;
  apparentTemperature: number;
  /** % */
  precipitationProbability: number;
  /** mm */
  precipitation: number;
  weatherCode: number;
  windSpeed: number;
  windDirection: number;
  isDay: boolean;
}

export interface DailyForecast {
  /** YYYY-MM-DD on the location's calendar */
  date: string;
  temperatureMin: number;
  temperatureMax: number;
  weatherCode: number;
  precipitationSum: number;
  precipitationProbabilityMax: number;
  windSpeedMax: number;
  windGustsMax: number;
  windDirectionDominant: number;
  uvIndexMax: number;
  sunrise: string;
  sunset: string;
  daylightSeconds: number;
}

export interface Forecast {
  current: CurrentForecast;
  /** From the current hour to the end of the last day. */
  hourly: HourlyForecast[];
  /** Seven days from today. */
  daily: DailyForecast[];
  /** The location's IANA zone: every time on the page is shown in it. */
  timezone: string;
  /** Metres */
  elevation: number;
  /** When the backend fetched it; much older than 10 minutes means Open-Meteo is failing. */
  fetchedAt: string;
}
