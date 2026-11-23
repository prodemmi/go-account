### Phase 1: Security & Core Logic (Highest Priority)
*Focus: Fixing vulnerabilities and refining authentication flows.*

**Prompt 1: Security Hardening (OTP & Routing)**
> **Objective:** Fix security vulnerabilities in the authentication flow.
> 1. **OTP Vulnerability:** In the `verify-otp` page, the system currently relies on OTP resource data extracted from the URL. This is insecure as it allows manual manipulation. Refactor the verification logic to ensure the OTP session/resource is retrieved from a secure, server-side state or a signed, encrypted token that cannot be forged via URL parameters.
> 2. **Auth Routing:** Implement a middleware/routing rule: When `custom_auth_enabled` is `true`, all requests to `/auth/**` must return a `404 Not Found` to prevent unauthorized access to default Aman pages.
> 3. **Session/Status Consistency:** Investigate the discrepancy between the `Session Table` and the `Statistics/Filter` logic. Ensure that status calculations (e.g., `expired`, `active`) are synchronized between the database query layer and the UI components. Fix the bug where filtering for "expired" sessions returns zero results despite data existing.

**Prompt 2: Authentication Architecture (CAS Standard)**
> **Objective:** Implement a flexible Authentication Service pattern (CAS-like) that supports both Managed and Custom Auth.
> 1. **Logic Branching:** Implement two distinct authentication flows:
>    - **Mode A (Managed):** Use default Aman authentication pages and global domain cookies.
>    - **Mode B (Custom):** Business implements their own UI. The business service must connect to the Aman API to validate credentials and receive a valid session cookie.
> 2. **Standardization:** Define the standard protocol for how the Business Service communicates with the Aman Service for credential validation and cookie issuance.
> 3. **Developer Documentation:** Create comprehensive guides for both Mode A and Mode B, including API specifications for the `check-login` endpoint and cookie management.

---

### Phase 2: Backend & Business Logic
*Focus: Expanding capabilities and data enrichment.*

**Prompt 3: Business Configuration & Geo-Fencing**
> **Objective:** Add granular control for Business owners regarding authentication.
> 1. **Geo-Restriction:** Implement a feature allowing businesses to select a whitelist of allowed countries for authentication.
> 2. **Captcha Extensibility:** Refactor the Captcha module to support multiple providers. Implement a provider interface. Currently, `ArCaptcha` is implemented; add `Google reCAPTCHA` and other industry-standard providers. Allow businesses to select their preferred provider via configuration.
> 3. **User Identity Enrichment:** Enhance the session/user data aggregation. If a user's precise location is unavailable, implement an IP-to-Location lookup to populate the device/location data.

---

### Phase 3: UI/UX & Frontend Refinement
*Focus: Component updates and design improvements.*

**Prompt 4: Global UI Improvements**
> **Objective:** Standardize UI components and improve usability.
> 1. **Drawer Implementation:** Move all "Detail" views from separate pages/modals to a standardized Left-side Drawer component.
> 2. **Data Copy Utility:** Implement a "Copy to Clipboard" icon/button next to all unique IDs in tables and detail panels globally.
> 3. **Device Detection:** In the device table, implement User-Agent parsing to show human-readable device names (e.g., "iPhone", "Chrome on Windows") instead of raw strings.
> 4. **Persian Localization:** Add a Persian (Jalali) date-range picker to all table filters and the global dashboard statistics. Add "Default: All Time" filter buttons.

**Prompt 5: Auth Page Customization & UX**
> **Objective:** Revamp the Custom Auth experience.
> 1. **Visual Design:** Redesign the customizable auth page layouts to be modern and professional. Improve the CSS architecture to allow businesses to customize text colors and layout structures.
> 2. **Dynamic OTP UX:** In the `verify-otp` flow:
>    - If a user changes their mobile/email and the input matches the previous one, display the remaining time for the *original* OTP.
>    - If the input is new, trigger the logic for a new OTP and display the new expiration timer.
> 3. **Context-Aware UI:** Ensure the "Edit mobile/email" text/labels adjust dynamically based on the selected Business Channel.

---

### Phase 4: Dashboard & Analytics
*Focus: Data visualization.*

**Prompt 6: Dashboard & Data Visualization**
> **Objective:** Enhance the analytical capabilities of the dashboard.
> 1. **Expanded Metrics:** Add new statistical charts and data visualizations to the dashboard.
> 2. **Integration:** Ensure all dashboard statistics are reactive to the new Persian date-range filters and global status filters.
> 3. **Data Accuracy:** Verify that all "Stats" cards correctly reflect the underlying data in the session tables (re-verify the Fix from Phase 1).

### Phase 5: Documentation Integration
*Focus: Improving user onboarding via contextual help.*

**Prompt 7: Contextual Documentation Links**
> **Objective:** Implement "Go-to-Guide" links within the Business Settings.
> 1. **Contextual Mapping:** For every configuration setting or input field within the Business Settings panel, add a subtle "Guide" icon or link.
> 2. **Deep Linking:** When clicked, the link must perform an anchor-jump or route the user to the specific corresponding section in the Developer Documentation/Guide.
> 3. **UX Flow:** Ensure the transition is seamless (e.g., opening the guide in a new tab or a side-panel) so the user does not lose their current progress in the settings form.
