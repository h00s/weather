export interface Forecast {
  latitude: number;
  longitude: number;
  generationTimeMs: number;
  utcOffsetSeconds: number;
  timezone: string;
  timezoneAbbreviation: string;
  elevation: number;
  current: {
    time: string;
    interval: number;
    data: {
      precipitation: number;
      relative_humidity_2m: number;
      temperature_2m: number;
      weather_code: number;
    };
  };
  minutely15: {
    time: null;
    data: null;
  };
  hourly: {
    time: string[];
    data: {
      precipitation: number[];
      relative_humidity_2m: number[];
      temperature_2m: number[];
      weather_code: number[];
    };
  };
  daily: {
    time: string[];
    data: {
      precipitation_sum: number[];
      temperature_2m_max: number[];
      temperature_2m_min: number[];
    };
  };
}