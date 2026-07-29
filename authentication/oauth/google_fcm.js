(() => {
	console.debug("[ls] google_fcm.js loaded");

	// Injected by the server per-response (see authentication/oauth/google_fcm.go).
	// Bound to the ls_oidc_fcm_nonce cookie set alongside this script. Google
	// embeds this value verbatim in the "nonce" claim of the ID token it
	// returns, and the server verifies that claim against the cookie before
	// trusting the credential. This prevents a credential obtained on one
	// browser from being replayed against a different victim's browser
	// (login CSRF).
	const SERVER_NONCE = "{{.Nonce}}";

	const currentPage =
		window.location.pathname + window.location.search + window.location.hash;

	const url = new URL(document.currentScript.src);
	const clientId = url.searchParams.get("client_id");

	if (!clientId) {
		throw new Error("[ls ; google_fcm.js] client_id is required");
	}

	const script = document.createElement("script");
	script.src = "https://accounts.google.com/gsi/client";
	script.async = true;
	script.defer = true;

	const hasCookie = (name) =>
		document.cookie
			.split("; ")
			.some((cookie) => cookie.startsWith(`${name}=`));

	script.onload = () => {
			google.accounts.id.initialize({
				client_id: clientId,
				nonce: SERVER_NONCE,
				callback: async (response) => {
					const form = document.createElement("form");
					form.method = "POST";
					form.action = `/api/login?provider=google&page=${encodeURIComponent(currentPage)}`;

					for (const [name, value] of Object.entries({
						credential: response.credential,
						select_by: response.select_by || "",
					})) {
						const input = document.createElement("input");
						input.type = "hidden";
						input.name = name;
						input.value = value;
						form.appendChild(input);
					}

					document.body.appendChild(form);
					form.submit();
				},
				auto_select: true,
				context: "use",
				color_scheme: "light",
				itp_support: false,
			});

			if (!hasCookie("ls_expires_at")) {
				google.accounts.id.prompt();
			} else {
				console.debug("[ls] google_fcm.js: ls_expires_at cookie found, skipping prompt")
			}
	};

	document.head.appendChild(script);
})();
