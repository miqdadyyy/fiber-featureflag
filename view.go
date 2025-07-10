package featureflag

const (
	IndexView = `<!DOCTYPE html>
<html lang='en'>
<head>
    <meta charset='UTF-8'>
    <meta name='viewport' content='width=device-width, initial-scale=1.0'>
    <title>Feature Flag Manager</title>
    <script src='https://cdn.tailwindcss.com'></script>
    <script src='https://unpkg.com/lucide@latest'></script>
    <link href='https://cdnjs.cloudflare.com/ajax/libs/toastr.js/latest/toastr.min.css' rel='stylesheet'>
    <style>
        body {
            font-family: 'Inter', sans-serif;
        }

        :root {
            --background: 0 0% 100%;
            --foreground: 222.2 84% 4.9%;
            --card: 0 0% 100%;
            --card-foreground: 222.2 84% 4.9%;
            --primary: 222.2 47.4% 11.2%;
            --primary-foreground: 210 20% 98%;
            --destructive: 0 84.2% 60.2%;
            --destructive-foreground: 210 20% 98%;
            --input: 214.3 31.8% 91.4%;
            --border: 214.3 31.8% 91.4%;
            --ring: 222.2 84% 4.9%;
            --accent: 210 40% 96.1%;
            --accent-foreground: 222.2 47.4% 11.2%;
            --muted: 210 40% 96.1%;
            --muted-foreground: 215.4 16.3% 46.9%;
        }

        .toggle-switch {
            --toggle-bg-off: hsl(var(--input));
            --toggle-bg-on: hsl(var(--primary));
            --toggle-handle-off: hsl(var(--background));
            --toggle-handle-on: hsl(var(--background));
            --toggle-border-off: hsl(var(--input));
            --toggle-border-on: hsl(var(--primary));
            --toggle-shadow: rgba(0, 0, 0, 0.1);
        }

        .toggle-switch input:checked + .slider {
            background-color: var(--toggle-bg-on);
            border-color: var(--toggle-border-on);
        }

        .toggle-switch input + .slider {
            background-color: var(--toggle-bg-off);
            border: 1px solid var(--toggle-border-off);
            transition: background-color 0.2s, border-color 0.2s;
        }

        .toggle-switch input:checked + .slider:before {
            transform: translateX(20px);
            background-color: var(--toggle-handle-on);
        }

        .toggle-switch input + .slider:before {
            content: '';
            position: absolute;
            height: 16px;
            width: 16px;
            left: 2px;
            bottom: 2px;
            background-color: var(--toggle-handle-off);
            border-radius: 9999px;
            transition: transform 0.2s, background-color 0.2s;
            box-shadow: 0 1px 2px var(--toggle-shadow);
        }

        input:focus-visible, button:focus-visible, textarea:focus-visible {
            outline: none;
            box-shadow: 0 0 0 2px hsl(var(--ring));
            outline-offset: 2px;
        }

        .btn {
            display: inline-flex;
            align-items: center;
            justify-content: center;
            white-space: nowrap;
            border-radius: 0.375rem;
            font-weight: 500;
            transition: background-color 0.2s, box-shadow 0.2s, opacity 0.2s;
            cursor: pointer;
            user-select: none;
            height: 2.5rem;
            padding-left: 1rem;
            padding-right: 1rem;
            font-size: 0.875rem;
            line-height: 1.25rem;
        }

        .btn-default {
            background-color: hsl(var(--primary));
            color: hsl(var(--primary-foreground));
            box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
        }

        .btn-default:hover {
            background-color: hsl(var(--primary) / 0.9);
            box-shadow: 0 2px 4px -1px rgba(0, 0, 0, 0.1), 0 1px 2px -1px rgba(0, 0, 0, 0.06);
        }

        .btn-outline {
            border: 1px solid hsl(var(--input));
            background-color: hsl(var(--background));
            color: hsl(var(--foreground));
        }

        .btn-outline:hover {
            background-color: hsl(var(--accent));
            color: hsl(var(--accent-foreground));
        }

        .btn-destructive {
            background-color: hsl(var(--destructive));
            color: hsl(var(--destructive-foreground));
        }

        .btn-destructive:hover {
            background-color: hsl(var(--destructive) / 0.9);
        }

        .input-field {
            display: flex;
            width: 100%;
            border: 1px solid hsl(var(--input));
            background-color: hsl(var(--background));
            color: hsl(var(--foreground));
            border-radius: 0.375rem;
            height: 2.5rem;
            padding-left: 0.75rem;
            padding-right: 0.75rem;
            font-size: 0.875rem;
            line-height: 1.25rem;
            transition: border-color 0.2s, box-shadow 0.2s;
        }

        .input-field:focus-visible {
            border-color: hsl(var(--ring));
            box-shadow: 0 0 0 2px hsl(var(--ring));
        }

        .textarea-field {
            display: flex;
            min-height: 5rem;
            width: 100%;
            border: 1px solid hsl(var(--input));
            background-color: hsl(var(--background));
            color: hsl(var(--foreground));
            border-radius: 0.375rem;
            padding: 0.75rem;
            font-size: 0.875rem;
            line-height: 1.25rem;
            transition: border-color 0.2s, box-shadow 0.2s;
            resize: vertical;
        }

        .textarea-field:focus-visible {
            border-color: hsl(var(--ring));
            box-shadow: 0 0 0 2px hsl(var(--ring));
        }

        .card {
            border-radius: 0.5rem;
            border: 1px solid hsl(var(--border));
            background-color: hsl(var(--card));
            color: hsl(var(--card-foreground));
            box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px 0 rgba(0, 0, 0, 0.06);
        }

        .card-header {
            display: flex;
            flex-direction: column;
            padding: 1.5rem;
        }

        .card-title {
            font-size: 1.5rem;
            font-weight: 600;
            line-height: 2rem;
        }

        .card-description {
            font-size: 0.875rem;
            color: hsl(var(--muted-foreground));
        }

        .card-content {
            padding: 1.5rem;
            padding-top: 0;
        }

        .badge {
            display: inline-flex;
            align-items: center;
            border-radius: 9999px;
            border: 1px solid transparent;
            padding-left: 0.5rem;
            padding-right: 0.5rem;
            font-size: 0.75rem;
            font-weight: 600;
            height: 1.25rem;
            text-transform: uppercase;
            transition: background-color 0.2s;
        }

        .badge-default {
            background-color: hsl(var(--primary));
            color: hsl(var(--primary-foreground));
        }

        .badge-secondary {
            background-color: hsl(var(--muted));
            color: hsl(var(--muted-foreground));
        }

        .dialog-overlay {
            position: fixed;
            inset: 0;
            background-color: rgba(0, 0, 0, 0.5);
            z-index: 50;
            opacity: 0;
            transition: opacity 0.3s ease-out;
            pointer-events: none;
        }

        .dialog-overlay.open {
            opacity: 1;
            pointer-events: auto;
        }

        .dialog-content {
            position: fixed;
            top: 50%;
            left: 50%;
            transform: translate(-50%, -48%) scale(0.95);
            background-color: hsl(var(--background));
            border-radius: 0.5rem;
            border: 1px solid hsl(var(--border));
            padding: 1.5rem;
            width: 90vw;
            max-width: 32rem;
            z-index: 50;
            opacity: 0;
            transition: opacity 0.3s ease-out, transform 0.3s ease-out;
            pointer-events: none;
        }

        .dialog-content.open {
            opacity: 1;
            transform: translate(-50%, -50%) scale(1);
            pointer-events: auto;
        }

        .dialog-header {
            display: flex;
            flex-direction: column;
            text-align: center;
            margin-bottom: 1rem;
        }

        .dialog-title {
            font-size: 1.125rem;
            font-weight: 600;
            line-height: 1.75rem;
        }

        .dialog-description {
            font-size: 0.875rem;
            color: hsl(var(--muted-foreground));
        }

        .dialog-footer {
            display: flex;
            justify-content: end;
            gap: 0.5rem;
            margin-top: 1.5rem;
        }

        .btn-icon {
            background-color: transparent;
            border: none;
            color: hsl(var(--muted-foreground));
            height: 2rem;
            width: 2rem;
            padding: 0;
            display: inline-flex;
            align-items: center;
            justify-content: center;
            border-radius: 0.375rem;
            transition: background-color 0.2s, color 0.2s;
        }

        .btn-icon:hover {
            background-color: hsl(var(--accent));
            color: hsl(var(--foreground));
        }
    </style>
</head>
<body class='min-h-screen flex flex-col items-center py-12 px-4 sm:px-6 lg:px-8' style='background-color: hsl(var(--background)); color: hsl(var(--foreground));'>
<div class='max-w-6xl mx-auto w-full'>
    <div class='mb-8'>
        <div class='flex items-center justify-between'>
            <div>
                <h1 class='text-3xl font-bold'>Feature Flags</h1>
                <p class='text-gray-600 mt-2'>Manage and control feature rollouts</p>
            </div>
            <button id='addFlagTrigger' class='btn btn-default'>
                <i data-lucide='plus' class='w-4 h-4 mr-2'></i>
                Add Flag
            </button>
        </div>
    </div>

    <div class='grid grid-cols-1 md:grid-cols-3 gap-6 mb-8'>
        <div class='card'>
            <div class='card-header flex flex-row items-center justify-between space-y-0 pb-2'>
                <h3 class='text-sm font-medium'>Total Flags</h3>
                <i data-lucide='settings' class='h-4 w-4' style='color: hsl(var(--muted-foreground));'></i>
            </div>
            <div class='card-content'>
                <div id='totalFlagsCount' class='text-2xl font-bold'>0</div>
                <p class='text-xs' style='color: hsl(var(--muted-foreground));'>Total feature flags</p>
            </div>
        </div>
        <div class='card'>
            <div class='card-header flex flex-row items-center justify-between space-y-0 pb-2'>
                <h3 class='text-sm font-medium'>Active Flags</h3>
                <i data-lucide='zap' class='h-4 w-4' style='color: hsl(var(--muted-foreground));'></i>
            </div>
            <div class='card-content'>
                <div id='activeFlagsCount' class='text-2xl font-bold'>0</div>
                <p class='text-xs' style='color: hsl(var(--muted-foreground));'>Currently enabled</p>
            </div>
        </div>
        <div class='card'>
            <div class='card-header flex flex-row items-center justify-between space-y-0 pb-2'>
                <h3 class='text-sm font-medium'>Coverage</h3>
                <i data-lucide='users' class='h-4 w-4' style='color: hsl(var(--muted-foreground));'></i>
            </div>
            <div class='card-content'>
                <div id='coveragePercentage' class='text-2xl font-bold'>0%</div>
                <p class='text-xs' style='color: hsl(var(--muted-foreground));'>Flags enabled</p>
            </div>
        </div>
    </div>

    <div class='card'>
        <div class='card-header'>
            <h2 class='card-title'>Feature Flags</h2>
            <p class='card-description'>
                Toggle feature flags to control feature availability across your application.
            </p>
        </div>
        <div class='card-content'>
            <div id='featureFlagsContainer' class='space-y-4'>
                <p id='loadingMessage' class='text-center text-lg' style='color: hsl(var(--muted-foreground));'>Loading feature flags...</p>
            </div>
        </div>
    </div>
</div>

<div id='addFlagDialogOverlay' class='dialog-overlay'></div>
<div id='addFlagDialogContent' class='dialog-content'>
    <div class='dialog-header'>
        <h3 class='dialog-title'>Create New Feature Flag</h3>
        <p class='dialog-description'>Add a new feature flag to control feature rollouts.</p>
    </div>
    <div class='grid gap-4 py-4'>
        <div class='grid gap-2'>
            <label for='newFlagNameInput' class='text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70'>Flag Name</label>
            <input
                    type='text'
                    id='newFlagNameInput'
                    placeholder='e.g., new-checkout-flow'
                    class='input-field'
            />
        </div>
    </div>
    <div class='dialog-footer'>
        <button id='cancelAddFlagButton' class='btn btn-outline'>
            Cancel
        </button>
        <button id='createFlagButton' class='btn btn-default'>
            Create Flag
        </button>
    </div>
</div>

<div id='deleteConfirmDialogOverlay' class='dialog-overlay'></div>
<div id='deleteConfirmDialogContent' class='dialog-content'>
    <div class='dialog-header'>
        <h3 class='dialog-title'>Confirm Deletion</h3>
        <p class='dialog-description'>Are you sure you want to delete the feature flag '<span id='flagNameToDelete' class='font-semibold'></span>'? This action cannot be undone.</p>
    </div>
    <div class='dialog-footer'>
        <button id='cancelDeleteButton' class='btn btn-outline'>
            Cancel
        </button>
        <button id='confirmDeleteButton' class='btn btn-destructive'>
            Delete
        </button>
    </div>
</div>

<script src='https://cdnjs.cloudflare.com/ajax/libs/jquery/3.6.0/jquery.min.js'></script>
<script src='https://cdnjs.cloudflare.com/ajax/libs/toastr.js/latest/toastr.min.js'></script>

<script>
  toastr.options = {
    'closeButton': true,
    'debug': false,
    'newestOnTop': true,
    'progressBar': true,
    'positionClass': 'toast-top-right',
    'preventDuplicates': false,
    'onclick': null,
    'showDuration': '300',
    'hideDuration': '1000',
    'timeOut': '5000',
    'extendedTimeOut': '1000',
    'showEasing': 'swing',
    'hideEasing': 'linear',
    'showMethod': 'fadeIn',
    'hideMethod': 'fadeOut'
  };

  const addFlagTrigger = document.getElementById('addFlagTrigger');
  const addFlagDialogOverlay = document.getElementById('addFlagDialogOverlay');
  const addFlagDialogContent = document.getElementById('addFlagDialogContent');
  const newFlagNameInput = document.getElementById('newFlagNameInput');
  const cancelAddFlagButton = document.getElementById('cancelAddFlagButton');
  const createFlagButton = document.getElementById('createFlagButton');

  const deleteConfirmDialogOverlay = document.getElementById('deleteConfirmDialogOverlay');
  const deleteConfirmDialogContent = document.getElementById('deleteConfirmDialogContent');
  const flagNameToDeleteSpan = document.getElementById('flagNameToDelete');
  const cancelDeleteButton = document.getElementById('cancelDeleteButton');
  const confirmDeleteButton = document.getElementById('confirmDeleteButton');

  const featureFlagsContainer = document.getElementById('featureFlagsContainer');
  const loadingMessage = document.getElementById('loadingMessage');
  const totalFlagsCount = document.getElementById('totalFlagsCount');
  const activeFlagsCount = document.getElementById('activeFlagsCount');
  const coveragePercentage = document.getElementById('coveragePercentage');

  let featureFlags = [];
  let flagToDeleteName = '';

  async function loadFeatureFlags() {
    try {
      const response = await fetch('', {
        method: 'GET',
        headers: {
          'Accept': 'application/json',
        },
      });
      
      if (response.ok) {
        const data = await response.json();
        featureFlags = Object.entries(data).map(([name, enabled]) => ({
          name,
          enabled,
        }));
      } else {
        console.error('Failed to load feature flags');
        featureFlags = [];
      }
    } catch (error) {
      console.error('Error loading feature flags:', error);
      featureFlags = [];
    }
  }

  function openDialog(overlayElement, contentElement) {
    overlayElement.classList.add('open');
    contentElement.classList.add('open');
  }

  function closeDialog(overlayElement, contentElement) {
    overlayElement.classList.remove('open');
    contentElement.classList.remove('open');
  }

  async function toggleFlag(name) {
    try {
      const response = await fetch('', {
        method: 'PATCH',
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ key: name }),
      });

      if (response.ok) {
        const result = await response.json();
        const flagIndex = featureFlags.findIndex(flag => flag.name === name);
        if (flagIndex !== -1) {
          featureFlags[flagIndex].enabled = result.enabled;
          renderFeatureFlags();
          toastr.success(` + "`Feature flag '${name}' is now ${result.enabled ? 'enabled' : 'disabled'}.`" + `);
        }
      } else {
        const error = await response.json();
        toastr.error(` + "`Failed to toggle feature flag: ${error.error}`" + `);
      }
    } catch (error) {
      console.error('Error toggling feature flag:', error);
      toastr.error('Failed to toggle feature flag');
    }
  }

  async function addFlag() {
    const name = newFlagNameInput.value.trim();

    if (!name) {
      toastr.error('Please fill in the flag name.');
      return;
    }

    if (featureFlags.some(flag => flag.name.toLowerCase() === name.toLowerCase())) {
      toastr.warning(` + "`Feature flag '${name}' already exists.`" + `);
      return;
    }

    try {
      const response = await fetch('', {
        method: 'POST',
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ key: name }),
      });

      if (response.ok) {
        const result = await response.json();
        const newFlag = {
          name: name,
          enabled: true,
        };
        featureFlags.push(newFlag);
        closeDialog(addFlagDialogOverlay, addFlagDialogContent);
        newFlagNameInput.value = '';
        renderFeatureFlags();
        toastr.success(` + "`Feature flag '${newFlag.name}' added successfully.`" + `);
      } else {
        const error = await response.json();
        toastr.error(` + "`Failed to add feature flag: ${error.error}`" + `);
      }
    } catch (error) {
      console.error('Error adding feature flag:', error);
      toastr.error('Failed to add feature flag');
    }
  }

  function openDeleteConfirmDialog(name) {
    flagToDeleteName = name;
    flagNameToDeleteSpan.textContent = name;
    openDialog(deleteConfirmDialogOverlay, deleteConfirmDialogContent);
  }

  async function confirmAndDeleteFlag() {
    try {
      const response = await fetch('', {
        method: 'DELETE',
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ key: flagToDeleteName }),
      });

      if (response.ok) {
        featureFlags = featureFlags.filter(flag => flag.name !== flagToDeleteName);
        renderFeatureFlags();
        toastr.info(` + "`Feature flag '${flagToDeleteName}' deleted.`" + `);
        closeDialog(deleteConfirmDialogOverlay, deleteConfirmDialogContent);
        flagToDeleteName = '';
      } else {
        const error = await response.json();
        toastr.error(` + "`Failed to delete feature flag: ${error.error}`" + `);
      }
    } catch (error) {
      console.error('Error deleting feature flag:', error);
      toastr.error('Failed to delete feature flag');
    }
  }

  function renderFeatureFlags() {
    featureFlagsContainer.innerHTML = '';
    loadingMessage.classList.add('hidden');

    if (featureFlags.length === 0) {
      featureFlagsContainer.innerHTML = "<p class='text-center text-lg py-4' style='color: hsl(var(--muted-foreground));'>No feature flags added yet.</p>";
    }

    let enabledCount = 0;

    featureFlags.forEach(flag => {
      if (flag.enabled) {
        enabledCount++;
      }

      const flagItem = document.createElement('div');
      flagItem.className = 'flex items-center justify-between p-4 border rounded-lg';
      flagItem.style.borderColor = 'hsl(var(--border))';
      flagItem.style.backgroundColor = 'hsl(var(--card))';

      flagItem.innerHTML = ` + "`<div class='flex-1'><div class='flex items-center gap-3 mb-2'><h3 class='font-semibold text-lg'>${flag.name}</h3><span class='badge ${flag.enabled ? 'badge-default' : 'badge-secondary'}'>${flag.enabled ? 'Enabled' : 'Disabled'}</span></div></div><div class='flex items-center space-x-2'><label for='flag-${flag.name}' class='text-sm font-medium'>${flag.enabled ? 'On' : 'Off'}</label><label class='toggle-switch relative inline-flex items-center cursor-pointer w-10 h-6'><input type='checkbox' id='flag-${flag.name}' class='sr-only peer' ${flag.enabled ? 'checked' : ''} data-flag-name='${flag.name}'><div class='slider absolute top-0 left-0 right-0 bottom-0 rounded-full'></div></label><button class='btn-icon delete-flag-btn' data-flag-name='${flag.name}'><i data-lucide='trash-2' class='w-4 h-4'></i></button></div>`" + `;
      featureFlagsContainer.appendChild(flagItem);

      const switchInput = flagItem.querySelector(` + "`#flag-${flag.name}`" + `);
      if (switchInput) {
        switchInput.addEventListener('change', (event) => {
          toggleFlag(event.target.dataset.flagName);
        });
      }

      const deleteButton = flagItem.querySelector(` + "`.delete-flag-btn[data-flag-name='${flag.name}']`" + `);
      if (deleteButton) {
        deleteButton.addEventListener('click', (event) => {
          openDeleteConfirmDialog(event.currentTarget.dataset.flagName);
        });
      }
    });

    totalFlagsCount.textContent = featureFlags.length;
    activeFlagsCount.textContent = enabledCount;
    coveragePercentage.textContent = featureFlags.length > 0 ? ` + "`${Math.round((enabledCount / featureFlags.length) * 100)}%`" + ` : '0%';

    lucide.createIcons();
  }

  addFlagTrigger.addEventListener('click', () => openDialog(addFlagDialogOverlay, addFlagDialogContent));
  cancelAddFlagButton.addEventListener('click', () => closeDialog(addFlagDialogOverlay, addFlagDialogContent));
  addFlagDialogOverlay.addEventListener('click', (event) => {
    if (event.target === addFlagDialogOverlay) {
      closeDialog(addFlagDialogOverlay, addFlagDialogContent);
    }
  });
  createFlagButton.addEventListener('click', addFlag);

  cancelDeleteButton.addEventListener('click', () => closeDialog(deleteConfirmDialogOverlay, deleteConfirmDialogContent));
  confirmDeleteButton.addEventListener('click', confirmAndDeleteFlag);
  deleteConfirmDialogOverlay.addEventListener('click', (event) => {
    if (event.target === deleteConfirmDialogOverlay) {
      closeDialog(deleteConfirmDialogOverlay, deleteConfirmDialogContent);
    }
  });

  document.addEventListener('DOMContentLoaded', async () => {
    await loadFeatureFlags();
    renderFeatureFlags();
  });
</script>
</body>
</html>`
)
