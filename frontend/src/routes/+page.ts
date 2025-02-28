import type { PageLoad } from './$types';
import { fetchForecast } from '$svc/forecast';
import type { Forecast } from '$type/Forecast';

export const load = (async ({ fetch }) => {
  const forecast = await fetchForecast(fetch);
  
  return {
    forecast
  };
}) satisfies PageLoad;