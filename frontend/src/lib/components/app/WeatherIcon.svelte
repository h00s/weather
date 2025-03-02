<script lang="ts">
  let { code } = $props();

  const size: number = 56;
  const color: string = "white";
  const secondaryColor: string = "#FFCB05";
  const strokeWidth: number = 2;

  function getWeatherSvg(weatherCode: number): string {
    switch (weatherCode) {
      // Clear sky
      case 0:
        return `
          <circle cx="50" cy="50" r="20" fill="${secondaryColor}" stroke="${secondaryColor}" />
          <line x1="50" y1="20" x2="50" y2="10" stroke="${secondaryColor}" />
          <line x1="50" y1="80" x2="50" y2="90" stroke="${secondaryColor}" />
          <line x1="20" y1="50" x2="10" y2="50" stroke="${secondaryColor}" />
          <line x1="80" y1="50" x2="90" y2="50" stroke="${secondaryColor}" />
          <line x1="25" y1="25" x2="18" y2="18" stroke="${secondaryColor}" />
          <line x1="75" y1="25" x2="82" y2="18" stroke="${secondaryColor}" />
          <line x1="25" y1="75" x2="18" y2="82" stroke="${secondaryColor}" />
          <line x1="75" y1="75" x2="82" y2="82" stroke="${secondaryColor}" />
        `;

      // Mainly clear
      case 1:
        return `
          <!-- Sun -->
          <circle cx="35" cy="40" r="15" fill="${secondaryColor}" stroke="${secondaryColor}" />
          <line x1="35" y1="20" x2="35" y2="15" stroke="${secondaryColor}" />
          <line x1="35" y1="60" x2="35" y2="65" stroke="${secondaryColor}" />
          <line x1="15" y1="40" x2="10" y2="40" stroke="${secondaryColor}" />
          <line x1="55" y1="40" x2="60" y2="40" stroke="${secondaryColor}" />
          
          <!-- Small Cloud -->
          <path
            d="M50,50 
              C50,45 60,45 60,50 
              C65,50 70,55 70,60 
              C70,65 65,70 60,70 
              C60,70 50,70 50,70"
            fill="white"
            stroke="${color}"
          />
        `;

      // Partly cloudy
      case 2:
        return `
          <!-- Sun -->
          <circle cx="35" cy="40" r="15" fill="${secondaryColor}" stroke="${secondaryColor}" />
          <line x1="35" y1="20" x2="35" y2="15" stroke="${secondaryColor}" />
          <line x1="35" y1="60" x2="35" y2="65" stroke="${secondaryColor}" />
          <line x1="15" y1="40" x2="10" y2="40" stroke="${secondaryColor}" />
          <line x1="55" y1="40" x2="60" y2="40" stroke="${secondaryColor}" />
          
          <!-- Medium Cloud -->
          <path
            d="M40,55 
              C40,45 60,45 60,55 
              C70,55 80,60 80,70 
              C80,80 65,85 55,85 
              C55,85 40,85 40,80"
            fill="white"
            stroke="${color}"
          />
        `;

      // Overcast
      case 3:
        return `
          <!-- Full Cloud -->
          <path
            d="M25,40 
              C25,30 45,30 45,40 
              C55,40 65,45 65,55 
              C65,65 55,70 45,70 
              C45,75 25,75 25,70
              C15,70 5,65 5,55
              C5,45 15,40 25,40"
            fill="white"
            stroke="${color}"
          />
        `;

      // Fog
      case 45:
        return `
          <line x1="20" y1="30" x2="80" y2="30" stroke="${color}" />
          <line x1="15" y1="40" x2="85" y2="40" stroke="${color}" />
          <line x1="20" y1="50" x2="80" y2="50" stroke="${color}" />
          <line x1="25" y1="60" x2="75" y2="60" stroke="${color}" />
          <line x1="30" y1="70" x2="70" y2="70" stroke="${color}" />
        `;

      // Depositing rime fog
      case 48:
        return `
          <line x1="20" y1="30" x2="80" y2="30" stroke="${color}" />
          <line x1="15" y1="40" x2="85" y2="40" stroke="${color}" />
          <line x1="20" y1="50" x2="80" y2="50" stroke="${color}" />
          <line x1="25" y1="60" x2="75" y2="60" stroke="${color}" />
          <line x1="30" y1="70" x2="70" y2="70" stroke="${color}" />
          <path d="M30,25 L25,15 M40,25 L45,15 M50,25 L48,15 M60,25 L63,15 M70,25 L68,15" stroke="${color}" />
        `;

      // Drizzle: Light
      case 51:
        return `
          <!-- Cloud -->
          <path
            d="M25,40 
              C25,30 45,30 45,40 
              C55,40 65,45 65,55 
              C65,65 55,70 45,70 
              C45,75 25,75 25,70
              C15,70 5,65 5,55
              C5,45 15,40 25,40"
            fill="white"
            stroke="${color}"
          />
          
          <!-- Light Drizzle -->
          <path d="M30,75 L28,82 M50,75 L48,82" stroke="${color}" />
        `;

      // Drizzle: Moderate
      case 53:
        return `
          <!-- Cloud -->
          <path
            d="M25,40 
              C25,30 45,30 45,40 
              C55,40 65,45 65,55 
              C65,65 55,70 45,70 
              C45,75 25,75 25,70
              C15,70 5,65 5,55
              C5,45 15,40 25,40"
            fill="white"
            stroke="${color}"
          />
          
          <!-- Moderate Drizzle -->
          <path d="M30,75 L28,82 M45,75 L43,82 M60,75 L58,82" stroke="${color}" />
        `;

      // Drizzle: Dense
      case 55:
        return `
          <!-- Cloud -->
          <path
            d="M25,40 
              C25,30 45,30 45,40 
              C55,40 65,45 65,55 
              C65,65 55,70 45,70 
              C45,75 25,75 25,70
              C15,70 5,65 5,55
              C5,45 15,40 25,40"
            fill="white"
            stroke="${color}"
          />
          
          <!-- Dense Drizzle -->
          <path d="M25,75 L23,82 M40,75 L38,82 M55,75 L53,82 M70,75 L68,82 M32,80 L30,87 M47,80 L45,87 M62,80 L60,87" stroke="${color}" />
        `;

      // Freezing Drizzle: Light
      case 56:
        return `
          <!-- Cloud -->
          <path
            d="M25,40 
              C25,30 45,30 45,40 
              C55,40 65,45 65,55 
              C65,65 55,70 45,70 
              C45,75 25,75 25,70
              C15,70 5,65 5,55
              C5,45 15,40 25,40"
            fill="white"
            stroke="${color}"
          />
          
          <!-- Light Freezing Drizzle -->
          <path d="M30,75 L28,82 M50,75 L48,82" stroke="${color}" />
          <path d="M28,82 L26,84 L28,86 M48,82 L46,84 L48,86" stroke="${color}" />
        `;

      // Freezing Drizzle: Dense
      case 57:
        return `
          <!-- Cloud -->
          <path
            d="M25,40 
              C25,30 45,30 45,40 
              C55,40 65,45 65,55 
              C65,65 55,70 45,70 
              C45,75 25,75 25,70
              C15,70 5,65 5,55
              C5,45 15,40 25,40"
            fill="white"
            stroke="${color}"
          />
          
          <!-- Dense Freezing Drizzle -->
          <path d="M30,75 L28,82 M45,75 L43,82 M60,75 L58,82" stroke="${color}" />
          <path d="M28,82 L26,84 L28,86 M43,82 L41,84 L43,86 M58,82 L56,84 L58,86" stroke="${color}" />
        `;

      // Rain: Slight
      case 61:
      case 80:
        return `
          <!-- Cloud -->
          <path
            d="M25,40 
              C25,30 45,30 45,40 
              C55,40 65,45 65,55 
              C65,65 55,70 45,70 
              C45,75 25,75 25,70
              C15,70 5,65 5,55
              C5,45 15,40 25,40"
            fill="white"
            stroke="${color}"
          />
          
          <!-- Slight Rain -->
          <path d="M30,75 L27,90 M50,75 L47,90" stroke="${color}" />
        `;

      // Rain: Moderate
      case 63:
      case 81:
        return `
          <!-- Cloud -->
          <path
            d="M25,40 
              C25,30 45,30 45,40 
              C55,40 65,45 65,55 
              C65,65 55,70 45,70 
              C45,75 25,75 25,70
              C15,70 5,65 5,55
              C5,45 15,40 25,40"
            fill="white"
            stroke="${color}"
          />
          
          <!-- Moderate Rain -->
          <path d="M30,75 L27,90 M45,75 L42,90 M60,75 L57,90" stroke="${color}" />
        `;

      // Rain: Heavy/Violent
      case 65:
      case 82:
        return `
          <!-- Cloud -->
          <path
            d="M25,40 
              C25,30 45,30 45,40 
              C55,40 65,45 65,55 
              C65,65 55,70 45,70 
              C45,75 25,75 25,70
              C15,70 5,65 5,55
              C5,45 15,40 25,40"
            fill="white"
            stroke="${color}"
          />
          
          <!-- Heavy/Violent Rain -->
          <path d="M25,75 L22,90 M40,75 L37,90 M55,75 L52,90 M70,75 L67,90" stroke="${color}" />
        `;

      // Freezing Rain: Light
      case 66:
        return `
          <!-- Cloud -->
          <path
            d="M25,40 
              C25,30 45,30 45,40 
              C55,40 65,45 65,55 
              C65,65 55,70 45,70 
              C45,75 25,75 25,70
              C15,70 5,65 5,55
              C5,45 15,40 25,40"
            fill="white"
            stroke="${color}"
          />
          
          <!-- Light Freezing Rain -->
          <path d="M30,75 L27,90 M50,75 L47,90" stroke="${color}" />
          <path d="M27,90 L24,93 L27,96 M47,90 L44,93 L47,96" stroke="${color}" />
        `;

      // Freezing Rain: Heavy
      case 67:
        return `
          <!-- Cloud -->
          <path
            d="M25,40 
              C25,30 45,30 45,40 
              C55,40 65,45 65,55 
              C65,65 55,70 45,70 
              C45,75 25,75 25,70
              C15,70 5,65 5,55
              C5,45 15,40 25,40"
            fill="white"
            stroke="${color}"
          />
          
          <!-- Heavy Freezing Rain -->
          <path d="M30,75 L27,90 M45,75 L42,90 M60,75 L57,90" stroke="${color}" />
          <path d="M27,90 L24,93 L27,96 M42,90 L39,93 L42,96 M57,90 L54,93 L57,96" stroke="${color}" />
        `;

      // Snow fall: Slight
      case 71:
      case 85:
        return `
          <!-- Cloud -->
          <path
            d="M25,40 
              C25,30 45,30 45,40 
              C55,40 65,45 65,55 
              C65,65 55,70 45,70 
              C45,75 25,75 25,70
              C15,70 5,65 5,55
              C5,45 15,40 25,40"
            fill="white"
            stroke="${color}"
          />
          
          <!-- Slight Snow -->
          <path d="M30,80 L30,90 M27,83 L33,87 M27,87 L33,83" stroke="${color}" />
          <path d="M50,80 L50,90 M47,83 L53,87 M47,87 L53,83" stroke="${color}" />
        `;

      // Snow fall: Moderate
      case 73:
        return `
          <!-- Cloud -->
          <path
            d="M25,40 
              C25,30 45,30 45,40 
              C55,40 65,45 65,55 
              C65,65 55,70 45,70 
              C45,75 25,75 25,70
              C15,70 5,65 5,55
              C5,45 15,40 25,40"
            fill="white"
            stroke="${color}"
          />
          
          <!-- Moderate Snow -->
          <path d="M30,80 L30,90 M27,83 L33,87 M27,87 L33,83" stroke="${color}" />
          <path d="M50,80 L50,90 M47,83 L53,87 M47,87 L53,83" stroke="${color}" />
          <path d="M70,80 L70,90 M67,83 L73,87 M67,87 L73,83" stroke="${color}" />
        `;

      // Snow fall: Heavy
      case 75:
      case 86:
        return `
          <!-- Cloud -->
          <path
            d="M25,40 
              C25,30 45,30 45,40 
              C55,40 65,45 65,55 
              C65,65 55,70 45,70 
              C45,75 25,75 25,70
              C15,70 5,65 5,55
              C5,45 15,40 25,40"
            fill="white"
            stroke="${color}"
          />
          
          <!-- Heavy Snow -->
          <path d="M20,80 L20,90 M17,83 L23,87 M17,87 L23,83" stroke="${color}" />
          <path d="M40,80 L40,90 M37,83 L43,87 M37,87 L43,83" stroke="${color}" />
          <path d="M60,80 L60,90 M57,83 L63,87 M57,87 L63,83" stroke="${color}" />
          <path d="M80,80 L80,90 M77,83 L83,87 M77,87 L83,83" stroke="${color}" />
        `;

      // Snow grains
      case 77:
        return `
          <!-- Cloud -->
          <path
            d="M25,40 
              C25,30 45,30 45,40 
              C55,40 65,45 65,55 
              C65,65 55,70 45,70 
              C45,75 25,75 25,70
              C15,70 5,65 5,55
              C5,45 15,40 25,40"
            fill="white"
            stroke="${color}"
          />
          
          <!-- Snow Grains (smaller snowflakes) -->
          <circle cx="30" cy="85" r="2" fill="white" stroke="${color}" />
          <circle cx="40" cy="80" r="2" fill="white" stroke="${color}" />
          <circle cx="50" cy="85" r="2" fill="white" stroke="${color}" />
          <circle cx="60" cy="80" r="2" fill="white" stroke="${color}" />
          <circle cx="70" cy="85" r="2" fill="white" stroke="${color}" />
        `;

      // Thunderstorm: Slight or moderate
      case 95:
        return `
          <!-- Cloud -->
          <path
            d="M25,40 
              C25,30 45,30 45,40 
              C55,40 65,45 65,55 
              C65,65 55,70 45,70 
              C45,75 25,75 25,70
              C15,70 5,65 5,55
              C5,45 15,40 25,40"
            fill="white"
            stroke="${color}"
          />
          
          <!-- Lightning -->
          <path d="M45,70 L40,85 L50,80 L35,100" fill="${secondaryColor}" stroke="${secondaryColor}" />
          
          <!-- Rain -->
          <path d="M30,75 L27,85 M60,75 L57,85" stroke="${color}" />
        `;

      // Thunderstorm with slight hail
      case 96:
        return `
          <!-- Cloud -->
          <path
            d="M25,40 
              C25,30 45,30 45,40 
              C55,40 65,45 65,55 
              C65,65 55,70 45,70 
              C45,75 25,75 25,70
              C15,70 5,65 5,55
              C5,45 15,40 25,40"
            fill="white"
            stroke="${color}"
          />
          
          <!-- Lightning -->
          <path d="M45,70 L40,85 L50,80 L35,100" fill="${secondaryColor}" stroke="${secondaryColor}" />
          
          <!-- Rain -->
          <path d="M30,75 L27,85 M60,75 L57,85" stroke="${color}" />
          
          <!-- Small Hail -->
          <circle cx="25" cy="90" r="3" fill="white" stroke="${color}" />
          <circle cx="65" cy="90" r="3" fill="white" stroke="${color}" />
          <circle cx="45" cy="95" r="3" fill="white" stroke="${color}" />
        `;

      // Thunderstorm with heavy hail
      case 99:
        return `
          <!-- Cloud -->
          <path
            d="M25,40 
              C25,30 45,30 45,40 
              C55,40 65,45 65,55 
              C65,65 55,70 45,70 
              C45,75 25,75 25,70
              C15,70 5,65 5,55
              C5,45 15,40 25,40"
            fill="white"
            stroke="${color}"
          />
          
          <!-- Lightning -->
          <path d="M45,70 L40,85 L50,80 L35,100" fill="${secondaryColor}" stroke="${secondaryColor}" />
          
          <!-- Rain -->
          <path d="M30,75 L27,85 M60,75 L57,85" stroke="${color}" />
          
          <!-- Large Hail -->
          <circle cx="25" cy="90" r="3" fill="white" stroke="${color}" />
          <circle cx="65" cy="90" r="3" fill="white" stroke="${color}" />
          <circle cx="45" cy="95" r="5" fill="white" stroke="${color}" />
        `;

      // Default (unknown weather code)
      default:
        return `
          <text x="50" y="60" text-anchor="middle" font-size="60" fill="${color}" stroke="none">?</text>
        `;
    }
  }
</script>

<svg
  xmlns="http://www.w3.org/2000/svg"
  width={size}
  height={size}
  viewBox="0 0 100 100"
  fill="none"
  stroke={color}
  stroke-width={strokeWidth}
  stroke-linecap="round"
  stroke-linejoin="round"
>
  {@html getWeatherSvg(code)}
</svg>
