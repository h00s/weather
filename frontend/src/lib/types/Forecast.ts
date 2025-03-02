export interface HourlyForecast {
  time: string[];
  data: {
    precipitation: number[];
    relative_humidity_2m: number[];
    temperature_2m: number[];
    weather_code: number[];
  };
}

export interface DailyForecast {
  time: string[];
  data: {
    precipitation_sum: number[];
    temperature_2m_max: number[];
    temperature_2m_min: number[];
    weather_code: number[];
  };
}

export interface CurrentForecast {
  time: string;
  interval: number;
  data: {
    precipitation: number;
    relative_humidity_2m: number;
    temperature_2m: number;
    weather_code: number;
  };
}

export interface Forecast {
  latitude: number;
  longitude: number;
  generationTimeMs: number;
  utcOffsetSeconds: number;
  timezone: string;
  timezoneAbbreviation: string;
  elevation: number;
  current: CurrentForecast;
  minutely15: {
    time: null;
    data: null;
  };
  hourly: HourlyForecast;
  daily: DailyForecast;
}