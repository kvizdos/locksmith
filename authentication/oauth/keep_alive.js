(function () {
	const getCookieValue = (name) =>
		document.cookie.match("(^|;)\\s*" + name + "\\s*=\\s*([^;]+)")?.pop() || "";

	const hasOAuthProvider = getCookieValue("ls_oauth_provider") !== "";

	window.onload = () => {
		if (hasOAuthProvider) {
			console.debug("[LS] Loaded OAuth Keep Alive");
			const iframe = document.createElement("iframe");
			iframe.src = "/api/auth/oauth/keep-alive";
			iframe.width = "0";
			iframe.height = "0";
			iframe.style.display = "none";
			document.body.appendChild(iframe);
		}
	};

	setInterval(() => {
		confirmGoodExpiration();
	}, 1000);

	/*
	While the middleware protects against initial page loads,
	if they are already on the page and their token expires, we
	need to handle that.
	*/
	let isRedirecting = false;
	const confirmGoodExpiration = () => {
		if (isRedirecting) {
			return;
		}
		if (+getCookieValue("ls_expires_at") < Date.now() / 1000) {
			if (!hasOAuthProvider) {
				console.debug("[LS] Redirecting to Login Page");

				window.location.href = "/login";
				return;
			}

			const fullPath = window.location.pathname + window.location.search;

			console.debug(
				`[LS] Redirecting to OAuth Provider (${getCookieValue("ls_oauth_provider")})`,
			);

			isRedirecting = true;
			// We're assuming here that the OAuth Keep Alive failed here,
			// force a full interactive one.
			window.location.href = `/api/auth/oauth/${getCookieValue("ls_oauth_provider")}?page=${fullPath}`;
		}
	};
})();
