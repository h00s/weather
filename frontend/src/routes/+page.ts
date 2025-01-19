import type { PageLoad } from './$types';
import { fetchForecast } from '$svc/forecast';

export const load: PageLoad = async ({ fetch }) => {
  const forecast = await fetchForecast(fetch);
  return { 
    forecast 
  };
};