(function () {
	const getCookieValue = (name) =>
		document.cookie.match("(^|;)\\s*" + name + "\\s*=\\s*([^;]+)")?.pop() || "";

	const provider = getCookieValue("ls_oauth_provider");
	const expiresAtCookie = getCookieValue("ls_expires_at");
	// ls_expires_at should be a UNIX timestamp (in seconds)
	const expiresAt = parseFloat(expiresAtCookie) || 0;
	const fullPath = window.location.pathname + window.location.search;

	console.debug("[LS] OAuth Provider:", provider);
	console.debug("[LS] Expires At (UNIX):", expiresAt);
	console.debug("[LS] Full Path:", fullPath);

	let alerted = false;
	let redirectTriggered = false;
	let runOnce = false;
	const checkExpiration = () => {
		const nowSec = Date.now() / 1000;
		const remaining = expiresAt - nowSec; // remaining time in seconds

		if (!runOnce) {
			console.debug(
				"[LS] Current Time (s):",
				nowSec,
				"Remaining (s):",
				remaining,
			);
			runOnce = true;
		}
		// If less than 5 minutes remain (or already expired), trigger the redirect immediately.
		if (remaining <= 300 && !redirectTriggered) {
			console.debug(
				"[LS] Less than 5 minutes remaining. Triggering redirect immediately.",
			);
			redirectTriggered = true;
			if (provider) {
				console.debug("[LS] Redirecting to OAuth keep alive endpoint.");
				window.location.href = `/api/auth/oauth/${provider}?page=${fullPath}`;
			} else {
				console.debug("[LS] No OAuth provider set. Redirecting to /login.");
				window.location.href = "/login";
			}
			return;
		}

		// If within 10 minutes (but more than 5 minutes remain), alert the user if not already alerted.
		if (remaining <= 600 && !alerted) {
			if (performance.now() <= 5000) {
				return;
			}
			alerted = true;
			const minutes = Math.floor(remaining / 60);
			const seconds = Math.floor(remaining % 60);
			if (provider) {
				console.debug(
					"[LS] Within 10 minutes of expiration with OAuth provider set. Alerting user.",
				);
				alert(
					`Your page will be refreshed in ${minutes - 5} minute${minutes - 5 !== 1 ? "s" : ""} and ${seconds} second${seconds !== 1 ? "s" : ""}. Please save any progress. We're sorry for the inconvenience.`,
				);
			} else {
				console.debug(
					"[LS] Within 10 minutes of expiration with no OAuth provider. Alerting user to save progress.",
				);
				alert(
					`Your session will expire in ${minutes - 5} minute${minutes - 5 !== 1 ? "s" : ""} and ${seconds} second${seconds !== 1 ? "s" : ""}. Please save any progress.`,
				);
			}
		}
	};

	console.debug("[LS] Initial expiration check.");
	// On initial load, if time remaining is less than 5 minutes, trigger the action immediately.
	checkExpiration();
	setInterval(checkExpiration, 1000);
})();
