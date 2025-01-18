import { PUBLIC_API_URL } from '$env/static/public';
import type { WeatherData } from '$lib/types/weather';

type FetchFunction = typeof fetch;

export async function fetchForecast(svelteFetch: FetchFunction = fetch): Promise<WeatherData> {
  return svelteFetch(`${PUBLIC_API_URL}/forecast?lat=45.5936&lon=17.2251`).then(res => res.json());
}