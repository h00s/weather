import { PUBLIC_API_URL } from '$env/static/public';
import type { Forecast } from '$type/Forecast';

type FetchFunction = typeof fetch;

export async function fetchForecast(
  svelteFetch: FetchFunction = fetch,
  lat = 45.5936,
  lon = 17.2251
): Promise<Forecast> {
  const url = new URL(`${PUBLIC_API_URL}/forecast`);
  url.searchParams.append('lat', lat.toString());
  url.searchParams.append('lon', lon.toString());
  
  const response = await svelteFetch(url.toString());
  
  if (!response.ok) {
    throw new Error(`API error: ${response.status}`);
  }
  
  return response.json();
}