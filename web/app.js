(() => {
  "use strict";

  const $ = (selector, root = document) => root.querySelector(selector);
  const $$ = (selector, root = document) => Array.from(root.querySelectorAll(selector));

  const els = {
    form: $("#qrForm"),
    bankField: $("#bankField"),
    bankTrigger: $("#bankTrigger"),
    bankTriggerLogo: $("#bankTriggerLogo"),
    bankTriggerInitial: $("#bankTriggerInitial"),
    bankTriggerImage: $("#bankTriggerImage"),
    bankTriggerText: $("#bankTriggerText"),
    bankError: $("#bankError"),
    bankDialog: $("#bankDialog"),
    bankSearch: $("#bankSearch"),
    bankSearchClear: $("#bankSearchClear"),
    bankList: $("#bankList"),
    bankEmpty: $("#bankEmpty"),
    account: $("#account"),
    accountClear: $("#accountClear"),
    accountError: $("#accountError"),
    accountName: $("#accountName"),
    accountNameClear: $("#accountNameClear"),
    accountNameCount: $("#accountNameCount"),
    accountNameError: $("#accountNameError"),
    amountTrigger: $("#amountTrigger"),
    amountField: $("#amountField"),
    amount: $("#amount"),
    amountClear: $("#amountClear"),
    amountError: $("#amountError"),
    messageTrigger: $("#messageTrigger"),
    messageField: $("#messageField"),
    message: $("#message"),
    messageClear: $("#messageClear"),
    messageCount: $("#messageCount"),
    messageError: $("#messageError"),
    generateBtn: $("#generateBtn"),
    generateLabel: $("#generateLabel"),
    resetBtn: $("#resetBtn"),
    stepInfo: $("#stepInfo"),
    stepDesign: $("#stepDesign"),
    resultEmpty: $("#resultEmpty"),
    resultFilled: $("#resultFilled"),
    resultCard: $(".result-card"),
    qrTicket: $("#qrTicket"),
    ticketQr: $("#ticketQr"),
    qrImage: $("#qrImage"),
    ticketAmount: $("#ticketAmount"),
    ticketMemo: $("#ticketMemo"),
    ticketBankLogo: $("#ticketBankLogo"),
    ticketBankInitial: $("#ticketBankInitial"),
    ticketName: $("#ticketName"),
    ticketAccount: $("#ticketAccount"),
    ticketBankName: $("#ticketBankName"),
    designBtn: $("#designBtn"),
    designDialog: $("#designDialog"),
    themeOptions: $$(".theme-option"),
    toggleLogos: $("#toggleLogos"),
    downloadBtn: $("#downloadBtn"),
    downloadMenu: $("#downloadMenu"),
    printBtn: $("#printBtn"),
    menuBtn: $("#menuBtn"),
    menuDialog: $("#menuDialog"),
    serverVersion: $("#serverVersion"),
    toast: $("#toast"),
    dialogs: $$(".dialog"),
  };

  const state = {
    banks: [],
    selectedBank: null,
    result: null,
    loading: false,
    abortController: null,
    requestSeq: 0,
    toastTimer: 0,
    theme: "light",
    showLogos: true,
  };

  const amountFormat = new Intl.NumberFormat("en-US", { maximumFractionDigits: 0 });
  const bankCollator = new Intl.Collator("vi", { sensitivity: "base", numeric: true });
  const POPULAR_CODES = [
    "MBBANK",
    "VIETCOMBANK",
    "BIDV",
    "VIETINBANK",
    "AGRIBANK",
    "ACB",
    "TECHCOMBANK",
    "TPBANK",
    "VPBANK",
    "SACOMBANK",
    "HDBANK",
    "SHB",
    "OCB",
    "VIB",
    "EXIMBANK",
    "MSB",
    "SEA_BANK",
    "PVCOM_BANK",
    "LIENVIETPOST_BANK",
  ];
  const CHECK_ICON =
    '<svg class="icon icon-xs" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M20 6 9 17l-5-5"/></svg>';
  const STAR_ICON =
    '<svg class="icon" viewBox="0 0 24 24" fill="#fbbf24" stroke="#f59e0b" stroke-width="1.4" stroke-linejoin="round" aria-hidden="true"><path d="m12 2.6 2.95 5.98 6.6 1-4.78 4.65 1.13 6.57L12 17.7l-5.9 3.1 1.13-6.57-4.78-4.65 6.6-1z"/></svg>';

  const fieldMap = {
    bank: { field: () => els.bankField, error: () => els.bankError, control: () => els.bankTrigger },
    account: { field: () => els.account.closest(".field"), error: () => els.accountError, control: () => els.account },
    accountName: { field: () => els.accountName.closest(".field"), error: () => els.accountNameError, control: () => els.accountName },
    amount: { field: () => els.amountField, error: () => els.amountError, control: () => els.amount },
    message: { field: () => els.messageField, error: () => els.messageError, control: () => els.message },
  };

  function init() {
    setOrigins();
    wireEvents();
    renderBankTrigger();
    updateAccountNameCount();
    updateMessageCount();
    syncClearButtons();
    updateGenerateState();
    loadBanks();
    loadVersion();
  }

  function setOrigins() {
    const origin = window.location.origin;
    $$(".api-origin").forEach((el) => {
      el.textContent = origin;
    });
  }

  async function loadBanks() {
    try {
      const response = await fetch("/api/banks", { headers: { Accept: "application/json" } });
      if (!response.ok) throw new Error("bad status");
      const data = await response.json();
      state.banks = Array.isArray(data.banks) ? data.banks : [];
      els.bankTrigger.disabled = state.banks.length === 0;
      renderBankTrigger();
      if (els.bankDialog.open) renderBankList(els.bankSearch.value);
    } catch {
      els.bankTrigger.disabled = true;
      els.bankTriggerText.replaceChildren(placeholderEl("Không tải được danh sách ngân hàng"));
      showToast("Không tải được danh sách ngân hàng.", "error");
    }
    updateGenerateState();
  }

  async function loadVersion() {
    try {
      const response = await fetch("/healthz", { headers: { Accept: "application/json" } });
      if (!response.ok) return;
      const data = await response.json();
      if (data && data.version) {
        els.serverVersion.textContent = data.version === "dev" ? "dev" : `v${data.version}`;
      }
    } catch {
      /* giữ giá trị mặc định */
    }
  }

  function wireEvents() {
    els.form.addEventListener("submit", (event) => {
      event.preventDefault();
      generate();
    });

    els.resetBtn.addEventListener("click", resetForm);

    els.bankTrigger.addEventListener("click", openBankDialog);
    els.bankSearch.addEventListener("input", () => {
      renderBankList(els.bankSearch.value);
      els.bankSearchClear.hidden = els.bankSearch.value.length === 0;
    });
    els.bankSearchClear.addEventListener("click", () => {
      els.bankSearch.value = "";
      els.bankSearchClear.hidden = true;
      renderBankList("");
      els.bankSearch.focus();
    });
    els.bankSearch.addEventListener("keydown", (event) => {
      if (event.key !== "Enter") return;
      const first = els.bankList.querySelector(".bank-item");
      if (!first) return;
      event.preventDefault();
      first.click();
    });
    els.bankList.addEventListener("click", (event) => {
      const item = event.target.closest(".bank-item");
      if (!item) return;
      const bank = state.banks.find((candidate) => candidate.code === item.dataset.code);
      if (bank) selectBank(bank);
    });
    els.bankList.addEventListener("keydown", (event) => {
      if (!["ArrowDown", "ArrowUp", "Home", "End"].includes(event.key)) return;
      const items = $$(".bank-item", els.bankList);
      if (items.length === 0) return;
      event.preventDefault();
      const current = items.indexOf(document.activeElement);
      let next = current;
      if (event.key === "Home") next = 0;
      else if (event.key === "End") next = items.length - 1;
      else if (event.key === "ArrowDown") next = current < 0 ? 0 : Math.min(current + 1, items.length - 1);
      else if (event.key === "ArrowUp") next = current <= 0 ? -1 : current - 1;
      if (next < 0) {
        els.bankSearch.focus();
        return;
      }
      items[next].focus();
    });
    els.bankDialog.addEventListener("close", () => {
      els.bankTrigger.setAttribute("aria-expanded", "false");
    });

    els.account.addEventListener("input", () => {
      const cleaned = els.account.value.replace(/\D/g, "").slice(0, 19);
      if (cleaned !== els.account.value) els.account.value = cleaned;
      clearFieldError("account");
      syncClearButtons();
      updateGenerateState();
    });

    els.accountName.addEventListener("input", () => {
      const cleaned = sanitizeAccountName(els.accountName.value);
      if (cleaned !== els.accountName.value) els.accountName.value = cleaned;
      clearFieldError("accountName");
      updateAccountNameCount();
      syncClearButtons();
      updateGenerateState();
    });

    els.amount.addEventListener("input", () => {
      formatAmountInput();
      clearFieldError("amount");
      syncClearButtons();
    });

    els.message.addEventListener("input", () => {
      const chars = Array.from(els.message.value);
      if (chars.length > 50) els.message.value = chars.slice(0, 50).join("");
      clearFieldError("message");
      updateMessageCount();
      syncClearButtons();
    });

    els.accountClear.addEventListener("click", () => {
      els.account.value = "";
      els.account.dispatchEvent(new Event("input", { bubbles: true }));
      els.account.focus();
    });
    els.accountNameClear.addEventListener("click", () => {
      els.accountName.value = "";
      els.accountName.dispatchEvent(new Event("input", { bubbles: true }));
      els.accountName.focus();
    });
    els.amountClear.addEventListener("click", () => {
      els.amount.value = "";
      els.amount.dispatchEvent(new Event("input", { bubbles: true }));
      els.amount.focus();
    });
    els.messageClear.addEventListener("click", () => {
      els.message.value = "";
      els.message.dispatchEvent(new Event("input", { bubbles: true }));
      els.message.focus();
    });

    els.amountTrigger.addEventListener("click", () => expandOptional("amount"));
    els.messageTrigger.addEventListener("click", () => expandOptional("message"));

    els.downloadBtn.addEventListener("click", () => {
      if (!state.result) return;
      if (els.downloadMenu.hidden) openDownloadMenu();
      else closeDownloadMenu({ restoreFocus: true });
    });
    els.downloadMenu.addEventListener("click", (event) => {
      const item = event.target.closest("button");
      if (!item) return;
      if (item.dataset.format) downloadResult(item.dataset.format);
      else if (item.dataset.action === "copy-link") copyResultLink();
      else if (item.dataset.action === "copy-payload") copyResultPayload();
    });
    els.printBtn.addEventListener("click", () => {
      if (!state.result) return;
      closeDownloadMenu();
      window.print();
    });
    els.designBtn.addEventListener("click", () => {
      if (!state.result) return;
      syncDesignControls();
      openDialog(els.designDialog);
    });

    els.themeOptions.forEach((option) => {
      option.addEventListener("click", () => {
        state.theme = option.dataset.theme;
        syncDesignControls();
      });
    });
    els.toggleLogos.addEventListener("change", () => {
      state.showLogos = els.toggleLogos.checked;
      applyTicketStyle();
    });

    els.menuBtn.addEventListener("click", () => openDialog(els.menuDialog));

    $$("[data-dialog]").forEach((button) => {
      button.addEventListener("click", () => openDialog(document.getElementById(button.dataset.dialog)));
    });
    $$("[data-close]").forEach((button) => {
      button.addEventListener("click", () => {
        const dialog = button.closest("dialog");
        if (dialog) dialog.close();
      });
    });
    $$("[data-toast]").forEach((button) => {
      button.addEventListener("click", () => showToast(button.dataset.toast));
    });
    els.dialogs.forEach((dialog) => {
      dialog.addEventListener("click", (event) => {
        if (event.target === dialog) dialog.close();
      });
    });

    document.addEventListener("pointerdown", (event) => {
      if (els.downloadMenu.hidden) return;
      if (!event.target.closest(".download-wrap")) closeDownloadMenu();
    });

    document.addEventListener("keydown", handleShortcut);
    els.qrImage.addEventListener("load", () => {
      if (!els.qrImage.getAttribute("src")) return;
      els.ticketQr.classList.remove("is-loading");
    });
    els.qrImage.addEventListener("error", () => {
      if (!els.qrImage.getAttribute("src")) return;
      els.ticketQr.classList.remove("is-loading");
      if (state.result) showToast("Không tải được ảnh QR.", "error");
    });
  }

  function handleShortcut(event) {
    if (event.ctrlKey || event.metaKey || event.altKey) return;

    if (event.key === "Escape" && !els.downloadMenu.hidden) {
      event.preventDefault();
      closeDownloadMenu({ restoreFocus: true });
      return;
    }

    if (document.querySelector("dialog[open]")) return;

    switch (event.key) {
      case "F1":
        event.preventDefault();
        openBankDialog();
        break;
      case "F2":
        event.preventDefault();
        els.account.focus();
        break;
      case "F3":
        event.preventDefault();
        els.accountName.focus();
        break;
      case "F4":
        event.preventDefault();
        expandOptional("amount");
        break;
      case "F5":
        event.preventDefault();
        expandOptional("message");
        break;
      default:
        break;
    }
  }

  function openBankDialog() {
    if (els.bankTrigger.disabled) return;
    clearFieldError("bank");
    els.bankSearch.value = "";
    els.bankSearchClear.hidden = true;
    renderBankList("");
    els.bankTrigger.setAttribute("aria-expanded", "true");
    if (!els.bankDialog.open) els.bankDialog.showModal();
    window.requestAnimationFrame(() => els.bankSearch.focus());
  }

  function renderBankList(query) {
    const normalized = normalizeSearch(query);
    els.bankList.replaceChildren();

    if (state.banks.length === 0) {
      els.bankEmpty.hidden = false;
      setBankEmptyText("Đang tải danh sách ngân hàng…");
      return;
    }

    if (normalized) {
      const matches = state.banks
        .filter((bank) => bankMatches(bank, normalized))
        .sort((a, b) => bankCollator.compare(a.shortName, b.shortName));
      if (matches.length === 0) {
        els.bankEmpty.hidden = false;
        setBankEmptyText("Không tìm thấy ngân hàng phù hợp.");
        return;
      }
      els.bankEmpty.hidden = true;
      els.bankList.append(bankSection("Kết quả tìm kiếm", matches, false));
      return;
    }

    els.bankEmpty.hidden = true;
    const popular = POPULAR_CODES.map((code) => state.banks.find((bank) => bank.code === code)).filter(Boolean);
    const popularCodes = new Set(popular.map((bank) => bank.code));
    const remaining = state.banks
      .filter((bank) => !popularCodes.has(bank.code))
      .sort((a, b) => bankCollator.compare(a.shortName, b.shortName));

    if (popular.length > 0) els.bankList.append(bankSection("Ngân hàng phổ biến", popular, true));
    if (remaining.length > 0) els.bankList.append(bankSection("Tất cả ngân hàng", remaining, false));
  }

  function setBankEmptyText(text) {
    const label = els.bankEmpty.querySelector("p");
    if (label) label.textContent = text;
  }

  function bankMatches(bank, normalized) {
    const haystack = normalizeSearch(
      [bank.shortName, bank.name, bank.code, bank.vietqrCode, bank.bin].filter(Boolean).join(" "),
    );
    return haystack.includes(normalized);
  }

  function bankSection(title, banks, withStar) {
    const section = document.createElement("section");
    section.className = "bank-section";
    const heading = document.createElement("h3");
    heading.className = "bank-section-title";
    if (withStar) heading.insertAdjacentHTML("beforeend", STAR_ICON);
    heading.append(document.createTextNode(title));
    const grid = document.createElement("div");
    grid.className = "bank-grid";
    banks.forEach((bank) => grid.append(bankButton(bank)));
    section.append(heading, grid);
    return section;
  }

  function bankButton(bank) {
    const button = document.createElement("button");
    button.type = "button";
    button.className = "bank-item";
    button.dataset.code = bank.code;
    button.title = bank.name;
    button.setAttribute("aria-pressed", String(state.selectedBank?.code === bank.code));
    if (state.selectedBank?.code === bank.code) button.classList.add("is-selected");

    const logo = document.createElement("span");
    logo.className = "bank-item-logo";
    logo.setAttribute("aria-hidden", "true");
    const initial = document.createElement("span");
    initial.className = "bank-initial";
    const image = document.createElement("img");
    image.alt = "";
    image.loading = "lazy";
    logo.append(initial, image);
    applyLogo(initial, image, bank);

    const copy = document.createElement("span");
    copy.className = "bank-item-copy";
    const name = document.createElement("span");
    name.className = "bank-item-name";
    name.textContent = bank.shortName;
    const sub = document.createElement("span");
    sub.className = "bank-item-sub";
    sub.textContent = bank.name;
    copy.append(name, sub);

    button.append(logo, copy);
    return button;
  }

  function selectBank(bank) {
    state.selectedBank = bank;
    clearFieldError("bank");
    renderBankTrigger();
    updateGenerateState();
    if (els.bankDialog.open) els.bankDialog.close();
    els.bankTrigger.focus();
  }

  function renderBankTrigger() {
    const bank = state.selectedBank;
    if (!bank) {
      els.bankTriggerLogo.hidden = true;
      els.bankTriggerImage.hidden = true;
      els.bankTriggerImage.removeAttribute("src");
      els.bankTriggerInitial.hidden = false;
      els.bankTriggerInitial.textContent = "";
      const text = state.banks.length > 0 ? "Chọn ngân hàng" : "Đang tải ngân hàng…";
      els.bankTriggerText.replaceChildren(placeholderEl(text));
      return;
    }
    els.bankTriggerLogo.hidden = false;
    applyLogo(els.bankTriggerInitial, els.bankTriggerImage, bank);
    const name = document.createElement("span");
    name.className = "bank-name";
    name.textContent = bank.shortName;
    const sub = document.createElement("span");
    sub.className = "bank-sub";
    sub.textContent = bank.name;
    els.bankTriggerText.replaceChildren(name, sub);
  }

  function placeholderEl(text) {
    const span = document.createElement("span");
    span.className = "bank-placeholder";
    span.textContent = text;
    return span;
  }

  function applyLogo(initial, image, bank, alt = "") {
    image.onerror = null;
    image.hidden = true;
    image.removeAttribute("src");
    image.alt = alt;
    initial.hidden = false;
    initial.textContent = initials(bank ? bank.shortName : "");
    if (!bank || !bank.logo) return;
    image.hidden = false;
    initial.hidden = true;
    image.onerror = () => {
      image.hidden = true;
      image.removeAttribute("src");
      initial.hidden = false;
    };
    image.src = bank.logo;
  }

  function expandOptional(kind) {
    const config = kind === "amount"
      ? { field: els.amountField, trigger: els.amountTrigger, input: els.amount }
      : { field: els.messageField, trigger: els.messageTrigger, input: els.message };
    if (!config.field.hidden) {
      config.input.focus();
      return;
    }
    config.field.hidden = false;
    config.trigger.hidden = true;
    config.trigger.setAttribute("aria-expanded", "true");
    window.requestAnimationFrame(() => config.input.focus());
  }

  function collapseOptional(kind) {
    const config = kind === "amount"
      ? { field: els.amountField, trigger: els.amountTrigger, input: els.amount }
      : { field: els.messageField, trigger: els.messageTrigger, input: els.message };
    config.field.hidden = true;
    config.trigger.hidden = false;
    config.trigger.setAttribute("aria-expanded", "false");
    config.input.value = "";
  }

  function sanitizeAccountName(value) {
    return value
      .normalize("NFD")
      .replace(/[\u0300-\u036f]/g, "")
      .replace(/đ/g, "D")
      .replace(/Đ/g, "D")
      .toUpperCase()
      .replace(/[^A-Z0-9 .,'&()\/-]/g, "")
      .replace(/\s+/g, " ")
      .slice(0, 50);
  }

  function formatAmountInput() {
    const digits = els.amount.value.replace(/\D/g, "").replace(/^0+(?=\d)/, "").slice(0, 12);
    els.amount.value = digits ? amountFormat.format(Number(digits)) : "";
  }

  function parseAmount() {
    const digits = els.amount.value.replace(/\D/g, "");
    return digits ? Number(digits) : 0;
  }

  function updateAccountNameCount() {
    els.accountNameCount.textContent = `(${els.accountName.value.length}/50 ký tự)`;
  }

  function updateMessageCount() {
    els.messageCount.textContent = `${Array.from(els.message.value).length}/50 ký tự`;
  }

  function syncClearButtons() {
    els.accountClear.hidden = els.account.value.length === 0;
    els.accountNameClear.hidden = els.accountName.value.length === 0;
    els.amountClear.hidden = els.amount.value.length === 0;
    els.messageClear.hidden = els.message.value.length === 0;
    els.bankSearchClear.hidden = els.bankSearch.value.length === 0;
  }

  function updateGenerateState() {
    els.generateBtn.disabled = state.loading;
  }

  function setGenerating(loading) {
    state.loading = loading;
    els.generateBtn.classList.toggle("is-loading", loading);
    els.generateBtn.setAttribute("aria-busy", String(loading));
    els.generateLabel.textContent = loading ? "Đang tạo…" : "Tạo mã QR";
    updateGenerateState();
  }

  function validateForm() {
    clearAllErrors();
    const invalid = [];

    if (!state.selectedBank) {
      setFieldError("bank", "Vui lòng chọn ngân hàng thụ hưởng.");
      invalid.push(els.bankTrigger);
    }

    const account = els.account.value.trim();
    if (!account) {
      setFieldError("account", "Vui lòng nhập số tài khoản.");
      invalid.push(els.account);
    } else if (account.length < 4 || account.length > 19) {
      setFieldError("account", "Số tài khoản phải từ 4 đến 19 ký tự.");
      invalid.push(els.account);
    }

    const name = els.accountName.value.trim();
    if (!name) {
      setFieldError("accountName", "Vui lòng nhập tên chủ tài khoản.");
      invalid.push(els.accountName);
    }

    if (Array.from(els.message.value).length > 50) {
      expandOptional("message");
      setFieldError("message", "Nội dung chuyển khoản tối đa 50 ký tự.");
      invalid.push(els.message);
    }

    if (invalid.length > 0) {
      const first = invalid[0];
      first.focus();
      if (first.scrollIntoView) first.scrollIntoView({ block: "center", behavior: "smooth" });
      return false;
    }
    return true;
  }

  async function generate() {
    if (state.loading) return;
    if (!validateForm()) return;

    const seq = ++state.requestSeq;
    if (state.abortController) state.abortController.abort();
    const controller = new AbortController();
    state.abortController = controller;
    setGenerating(true);
    closeDownloadMenu();

    try {
      const response = await fetch("/api/qr", {
        method: "POST",
        headers: { "Content-Type": "application/json", Accept: "application/json" },
        body: JSON.stringify({
          bank: state.selectedBank.code,
          account: els.account.value.trim(),
          accountName: els.accountName.value.trim(),
          amount: parseAmount(),
          message: els.message.value.trim(),
        }),
        signal: controller.signal,
      });
      const data = await response.json().catch(() => ({}));

      if (seq !== state.requestSeq) return;

      if (!response.ok) {
        clearResult();
        handleServerError(data);
        return;
      }

      renderResult(data);
      showToast("Đã tạo mã QR.");
      revealResult();
    } catch (error) {
      if (error.name === "AbortError" || seq !== state.requestSeq) return;
      clearResult();
      showToast("Không kết nối được máy chủ. Vui lòng thử lại.", "error");
    } finally {
      if (seq === state.requestSeq) {
        state.abortController = null;
        setGenerating(false);
      }
    }
  }

  function handleServerError(data) {
    const message = data && data.error ? data.error : "Không thể tạo mã QR.";
    if (data && data.field) {
      if (!fieldMap[data.field]) {
        showToast(message, "error");
        return;
      }
      if (data.field === "amount") expandOptional("amount");
      if (data.field === "message") expandOptional("message");
      setFieldError(data.field, message);
      const control = fieldMap[data.field].control();
      if (control) control.focus();
      return;
    }
    showToast(message, "error");
  }

  function renderResult(data) {
    state.result = data;
    const amount = Number(data.amount) || 0;
    const bank = data.bank || {};

    els.ticketAmount.textContent = amount > 0 ? `${amountFormat.format(amount)} VNĐ` : "Người nhận tự nhập";
    els.ticketMemo.hidden = !data.message;
    els.ticketMemo.textContent = data.message ? `ND: ${data.message}` : "";
    els.ticketName.hidden = !data.accountName;
    els.ticketName.textContent = data.accountName ? `Chủ TK: ${data.accountName}` : "";
    els.ticketAccount.textContent = `Số TK: ${data.account || ""}`;
    els.ticketBankName.textContent = bank.name || bank.shortName || "";
    applyLogo(els.ticketBankInitial, els.ticketBankLogo, bank, bank.shortName || "Ngân hàng");

    applyTicketStyle();
    els.resultEmpty.hidden = true;
    els.resultFilled.hidden = false;

    els.ticketQr.classList.add("is-loading");
    els.qrImage.alt = `Mã QR chuyển khoản ${bank.shortName || ""} · ${data.account || ""}`;
    els.qrImage.src = data.qr.png;
    if (els.qrImage.complete && els.qrImage.naturalWidth > 0) {
      els.ticketQr.classList.remove("is-loading");
    }

    els.downloadBtn.disabled = false;
    els.printBtn.disabled = false;
    setStepGenerated(true);
  }

  function applyTicketStyle() {
    els.qrTicket.classList.remove("theme-light", "theme-warm", "theme-dark");
    els.qrTicket.classList.add(`theme-${state.theme}`);
    els.qrTicket.classList.toggle("hide-logos", !state.showLogos);
  }

  function syncDesignControls() {
    els.themeOptions.forEach((option) => {
      const active = option.dataset.theme === state.theme;
      option.classList.toggle("is-active", active);
      option.setAttribute("aria-pressed", String(active));
    });
    els.toggleLogos.checked = state.showLogos;
    applyTicketStyle();
  }

  function setStepGenerated(generated) {
    const number = els.stepInfo.querySelector(".step-num");
    els.stepInfo.classList.toggle("is-active", !generated);
    els.stepInfo.classList.toggle("is-complete", generated);
    els.stepDesign.classList.toggle("is-active", generated);
    if (number) {
      if (generated) number.innerHTML = CHECK_ICON;
      else number.textContent = "1";
    }
  }

  function clearResult() {
    state.result = null;
    els.resultFilled.hidden = true;
    els.resultEmpty.hidden = false;
    els.ticketQr.classList.remove("is-loading");
    els.qrImage.removeAttribute("src");
    els.qrImage.alt = "";
    els.downloadBtn.disabled = true;
    els.printBtn.disabled = true;
    closeDownloadMenu();
    setStepGenerated(false);
  }

  function resetForm() {
    state.requestSeq += 1;
    if (state.abortController) {
      state.abortController.abort();
      state.abortController = null;
    }
    els.form.reset();
    state.selectedBank = null;
    state.loading = false;
    els.generateBtn.classList.remove("is-loading");
    els.generateBtn.removeAttribute("aria-busy");
    els.generateLabel.textContent = "Tạo mã QR";
    clearAllErrors();
    collapseOptional("amount");
    collapseOptional("message");
    updateAccountNameCount();
    updateMessageCount();
    syncClearButtons();
    renderBankTrigger();
    clearResult();
    updateGenerateState();
    if (!els.bankTrigger.disabled) els.bankTrigger.focus();
    else els.account.focus();
    showToast("Đã đặt lại thông tin.");
  }

  function revealResult() {
    if (!window.matchMedia("(max-width: 1024px)").matches) return;
    const reduceMotion = window.matchMedia("(prefers-reduced-motion: reduce)").matches;
    window.requestAnimationFrame(() => {
      els.resultCard.scrollIntoView({ behavior: reduceMotion ? "auto" : "smooth", block: "start" });
    });
  }

  function openDownloadMenu() {
    els.downloadMenu.hidden = false;
    els.downloadBtn.setAttribute("aria-expanded", "true");
    const first = els.downloadMenu.querySelector(".menu-item");
    if (first) first.focus();
  }

  function closeDownloadMenu(options = {}) {
    if (els.downloadMenu.hidden) return;
    els.downloadMenu.hidden = true;
    els.downloadBtn.setAttribute("aria-expanded", "false");
    if (options.restoreFocus) els.downloadBtn.focus();
  }

  function downloadResult(format) {
    if (!state.result) return;
    const href = format === "svg" ? state.result.qr.svg : state.result.qr.png;
    const anchor = document.createElement("a");
    anchor.href = href;
    anchor.download = `${resultSlug()}.${format === "svg" ? "svg" : "png"}`;
    document.body.appendChild(anchor);
    anchor.click();
    anchor.remove();
    closeDownloadMenu();
  }

  function copyResultLink() {
    if (!state.result) return;
    const url = new URL(state.result.qr.png, window.location.href).href;
    copyText(url).then((ok) => {
      showToast(ok ? "Đã sao chép liên kết ảnh." : "Không thể sao chép liên kết.", ok ? "success" : "error");
    });
    closeDownloadMenu();
  }

  function copyResultPayload() {
    if (!state.result) return;
    copyText(state.result.payload).then((ok) => {
      showToast(ok ? "Đã sao chép payload." : "Không thể sao chép payload.", ok ? "success" : "error");
    });
    closeDownloadMenu();
  }

  function resultSlug() {
    const bank = state.result.bank.vietqrCode || state.result.bank.code;
    const clean = (value) => String(value || "").replace(/[^0-9A-Za-z_-]/g, "");
    return `vietqr-${clean(bank)}-${clean(state.result.account)}`;
  }

  function openDialog(dialog) {
    if (!dialog || dialog.open) return;
    if (els.menuDialog.open) els.menuDialog.close();
    dialog.showModal();
  }

  function setFieldError(name, message) {
    const config = fieldMap[name];
    if (!config) return;
    const error = config.error();
    const field = config.field();
    const control = config.control();
    if (error) {
      error.textContent = message;
      error.hidden = false;
    }
    if (field) field.classList.add("is-invalid");
    if (control) control.setAttribute("aria-invalid", "true");
  }

  function clearFieldError(name) {
    const config = fieldMap[name];
    if (!config) return;
    const error = config.error();
    const field = config.field();
    const control = config.control();
    if (error) {
      error.textContent = "";
      error.hidden = true;
    }
    if (field) field.classList.remove("is-invalid");
    if (control) control.removeAttribute("aria-invalid");
  }

  function clearAllErrors() {
    Object.keys(fieldMap).forEach(clearFieldError);
  }

  function normalizeSearch(value) {
    return String(value || "")
      .normalize("NFD")
      .replace(/[\u0300-\u036f]/g, "")
      .replace(/đ/g, "d")
      .replace(/Đ/g, "D")
      .toLowerCase()
      .replace(/[^a-z0-9]/g, "");
  }

  function initials(name) {
    const words = String(name || "")
      .normalize("NFD")
      .replace(/[\u0300-\u036f]/g, "")
      .replace(/đ/gi, "d")
      .toUpperCase()
      .replace(/[^A-Z0-9\s]/g, " ")
      .trim()
      .split(/\s+/)
      .filter(Boolean);
    if (words.length === 0) return "?";
    if (words.length === 1) return words[0].slice(0, 2);
    return words[0][0] + words[words.length - 1][0];
  }

  async function copyText(text) {
    try {
      if (navigator.clipboard && window.isSecureContext) {
        await navigator.clipboard.writeText(text);
        return true;
      }
    } catch {
      /* thử fallback bên dưới */
    }

    try {
      const helper = document.createElement("textarea");
      helper.value = text;
      helper.setAttribute("readonly", "");
      helper.style.position = "fixed";
      helper.style.left = "-9999px";
      helper.style.top = "0";
      document.body.appendChild(helper);
      helper.select();
      const ok = document.execCommand("copy");
      helper.remove();
      return ok;
    } catch {
      return false;
    }
  }

  function showToast(message, kind = "success") {
    window.clearTimeout(state.toastTimer);
    els.toast.textContent = message;
    els.toast.dataset.variant = kind;
    els.toast.classList.add("is-visible");
    state.toastTimer = window.setTimeout(() => {
      els.toast.classList.remove("is-visible");
    }, 2600);
  }

  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", init, { once: true });
  } else {
    init();
  }
})();
