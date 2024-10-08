import type { CapacitorConfig } from '@capacitor/cli';

const config: CapacitorConfig = {
  appId: 'io.ionic.starter',
  appName: 'hotupdate',
  webDir: 'www/browser',
  plugins: {
    CapacitorUpdater: {
      autoUpdate: false,
	    statsUrl: ''
    }
  }
};

export default config;
