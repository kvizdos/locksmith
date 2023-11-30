class EphemeralTokenManager {
  static SharedInstance = new EphemeralTokenManager();

  constructor() {
    this.hasEphemeralToken =
      sessionStorage.getItem("eph-loaded") === true || false;
    this.loadingEphemeral = null;
  }

  async getToken() {
    if (this.loadingEphemeral != null) {
      await this.loadingEphemeral;
      return;
    }

    this.loadingEphemeral = new Promise(async (resolve) => {
      const fingerprint = await this.GenerateFingerprint();
      this.hasEphemeralToken = true;
      // ODO: Load Ephemeral Token Here
      setTimeout(() => {
        console.log("Done loading ephemeral!");
        resolve();
      }, 5000);
    });

    await this.loadingEphemeral;
  }

  async sha256(data) {
    const encoder = new TextEncoder();
    const dataBuffer = encoder.encode(data);
    const hashBuffer = await crypto.subtle.digest("SHA-256", dataBuffer);
    return Array.from(new Uint8Array(hashBuffer))
      .map((b) => b.toString(16).padStart(2, "0"))
      .join("");
  }

  async GenerateFingerprint() {
    const canvasFingerprint = async () => {
      const canvas = document.createElement("canvas");
      canvas.width = 500;
      canvas.height = 500;
      canvas.style.display = "none";
      document.body.appendChild(canvas);

      const ctx = canvas.getContext("2d");
      ctx.textBaseline = "top";
      ctx.font = "14px 'Arial'";
      ctx.textBaseline = "alphabetic";
      ctx.fillStyle = "#f60";
      ctx.fillRect(125, 1, 62, 20);
      ctx.fillStyle = "#069";
      ctx.fillText("Hello, world!", 2, 15);
      ctx.fillStyle = "rgba(102, 204, 0, 0.7)";
      ctx.fillText("Hello, world!", 4, 17);

      ctx.fillText("🤙", 100, 20);
      ctx.fillText("🎉", 110, 25);
      ctx.fillText("🤣", 115, 30);

      const dataUrl = canvas.toDataURL();

      canvas.remove();
      const hash = await this.sha256(dataUrl);
      return hash;
    };

    const getWebGLContext = () => {
      const canvas = document.createElement("canvas");
      let gl;
      try {
        gl =
          canvas.getContext("webgl") || canvas.getContext("experimental-webgl");
      } catch (e) {
        console.error("Failed to get WebGL context: ", e);
      }
      return gl;
    };

    const getWebGLRendererInfo = async () => {
      const gl = getWebGLContext();
      if (!gl) {
        return null;
      }

      const debugInfo = gl.getExtension("WEBGL_debug_renderer_info");
      if (debugInfo) {
        const info = {
          renderer: gl.getParameter(debugInfo.UNMASKED_RENDERER_WEBGL),
          vendor: gl.getParameter(debugInfo.UNMASKED_VENDOR_WEBGL),
        };

        const hash = await this.sha256(JSON.stringify(info));
        return hash;
      }

      return await this.sha256(JSON.stringify("blank"));
    };

    const hasTouchSupport = () => {
      return (
        "ontouchstart" in window ||
        navigator.maxTouchPoints > 0 ||
        navigator.msMaxTouchPoints > 0
      );
    };

    const getTouchPoints = () => {
      return navigator.maxTouchPoints || navigator.msMaxTouchPoints || 0;
    };

    const getTouchFingerprint = () => {
      return {
        touchSupport: hasTouchSupport(),
        maxTouchPoints: getTouchPoints(),
      };
    };

    const checkBatterySupport = async () => {
      if (navigator.getBattery) {
        try {
          const battery = await navigator.getBattery();
          return battery !== null;
        } catch {
          return false;
        }
      } else {
        return false;
      }
    };

    const generateAudioFingerprint = async () => {
      // Create an OfflineAudioContext
      const audioContext = new OfflineAudioContext(1, 44100, 44100);

      // Create an oscillator and connect it to the context
      const oscillator = audioContext.createOscillator();
      oscillator.type = "sine";
      oscillator.frequency.setValueAtTime(1000, audioContext.currentTime);
      oscillator.connect(audioContext.destination);
      oscillator.start(0);

      // Render the audio
      const audioBuffer = await audioContext.startRendering();
      const audioData = audioBuffer.getChannelData(0);

      // Convert the audio data to a Uint8Array
      const data = new Uint8Array(audioData.length);
      for (let i = 0; i < audioData.length; i++) {
        // Normalize the float32 audio data into uint8
        data[i] = Math.floor((audioData[i] * 0.5 + 0.5) * 255);
      }

      // Hash the data using SHA-256
      const hashBuffer = await crypto.subtle.digest("SHA-256", data);
      const hashArray = Array.from(new Uint8Array(hashBuffer));
      const hashHex = hashArray
        .map((b) => b.toString(16).padStart(2, "0"))
        .join("");

      return hashHex;
    };

    // Relatively Static Variables
    const screenHeight = window.screen.height * window.devicePixelRatio;
    const screenWidth = window.screen.width * window.devicePixelRatio;
    const colorDepth = window.screen.colorDepth;
    console.log("Color Depth", colorDepth);
    const screen = await this.sha256(
      JSON.stringify({
        screenHeight,
        screenWidth,
        colorDepth,
      }),
    );
    const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone;
    const hardwareConcurrency = window.navigator.hardwareConcurrency;
    const deviceMemory = navigator.deviceMemory || 0;
    const lang = window.navigator.language;
    const canvas = await canvasFingerprint();
    const webgl = await getWebGLRendererInfo();
    const touch = await this.sha256(JSON.stringify(getTouchFingerprint()));
    const battery = await checkBatterySupport();
    const platform =
      navigator?.userAgentData?.platform || navigator?.platform || "unknown";
    const audio = await generateAudioFingerprint();

    let fingerprint = {
      screen,
      timezone,
      hardwareConcurrency,
      deviceMemory,
      canvas,
      lang,
      webgl,
      touch,
      battery,
      platform,
      audio,
    };

    // Frequently Changing
    const getMediaDevices = async () => {
      try {
        const devices = await navigator.mediaDevices.enumerateDevices();
        return devices.map((device) => ({
          kind: device.kind,
          label: device.label,
          deviceId: device.deviceId,
          groupId: device.groupId,
        }));
      } catch (error) {
        return [];
      }
    };
    const userAgent = await this.sha256(navigator.userAgent);
    const windowSize = await this.sha256(
      JSON.stringify({
        height: window.innerHeight,
        width: window.innerWidth,
      }),
    );
    const doNotTrack = navigator.doNotTrack || false;
    const rawdevices = await getMediaDevices();
    const devices = await this.sha256(
      rawdevices.map((device) => Object.values(device).join(":")).join(";"),
    );

    fingerprint = {
      ...fingerprint,
      userAgent,
      windowSize,
      dnt: doNotTrack,
      devices,
    };

    return fingerprint;
  }
}

async function SecureFetch(url, options) {
  const hasEph = EphemeralTokenManager.SharedInstance.hasEphemeralToken;
  if (!hasEph) {
    await EphemeralTokenManager.SharedInstance.getToken();
  }

  return await fetch(url, options);
}

export { EphemeralTokenManager, SecureFetch };
