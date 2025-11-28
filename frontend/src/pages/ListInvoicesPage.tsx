import { useNavigate } from 'react-router-dom';

export default function ListInvoicesPage() {
  const navigate = useNavigate();

  const invoices = [
    {
      id: '2024-07-21A',
      title: 'Invoice #2024-07-21A',
      status: 'done',
      statusText: 'Extraction Done',
      thumbnail: 'https://lh3.googleusercontent.com/aida-public/AB6AXuBKeOFth3tWJ_-ufHJnhGVq2K26HfBpX5gFfOELUzQPD7KaDxOVf9wcSC75wCu3KU3sZ03T-E1rn2n0wNYfFtpA_G_3_E-OrXyBZDw6ryFXZdqaORz5JA3bs-az-lXLJk-UnYm6-7PDhS79zn1Yo1xFc8kd_XyP1LS8cgig6C1_e4hdJ08f91wAHjpd59wiGzL37kCT5KDoYrMTluKX6X2fuWdC2Cz_kfg13OAsdRvX0YJhKCBNBRLxDtLeclMsBZlJ2fzngWIxAmlc'
    },
    {
      id: 'costco',
      title: 'Costco Wholesale',
      status: 'review',
      statusText: 'Review Needed',
      thumbnail: 'https://lh3.googleusercontent.com/aida-public/AB6AXuDqIPAHrmQFxO73mNDELunGSNkzeXiZz-4q_n2VOP9u0RdbJGQHRmwv6THVMVmBy1NGCzM0n32R8gfEoLTzAJFm_IZANkK0iHNiEbGYXhVJTfxyqt6wKFUArT49CxMRVEUbTn9FduN1KvHYLbWSOi71gEKxZTsTBFaNXDpp38PbuFb2rN6y_6FQ_r5rpQdC3t1MOs_I4AkLxYRE_Y_am7OLtFgjZRdn7IW81DU9lmsqq7-SmqCmQC_3GRpDcNk3eO8LQIKlBgUhaJF3'
    },
    {
      id: 'office-depot',
      title: 'Office Depot',
      status: 'extracting',
      statusText: 'Extracting...',
      thumbnail: 'https://lh3.googleusercontent.com/aida-public/AB6AXuDdKGLhqQtsr-5L14Wv9HOug420m0Dw0m8xXZrYjGuQElh5ILX3O_6pdWh8xsOMoQ8tkychSaAG9tW8rI9tThAI3x6ZH3n9vkazUFlYzimgreThv6ffwJLMEFkh-YIc2jJcjYKhY6w-osv5k58mVfemiF9IpvtIMw7PdcXb3REoCDdx2r8uKJpX95th18mBdG-RITe66Ap8gwB6cUXP1B5oEIvC55sLexReh1M2Y5_C5AgdXrCWBtIyO_Ppp2F0YNfB_UqTMyWSub-U'
    },
    {
      id: 'utility',
      title: 'Utility Bill - May',
      status: 'done',
      statusText: 'Extraction Done',
      thumbnail: 'https://lh3.googleusercontent.com/aida-public/AB6AXuDorFqO0QbrCqR839MMP1c-48fxUhPmcHijsdKoHrGzK6DvKPKJ-w_sPkWFVp1vFAXnJF1UE8e_uIuTWZEaGX6hdDdsonm9VZ49j2J41W1eSEBK248aHuNJVHXQAkwQcZDMiU7Shpl8HQWIXIG3cUpFB0bLExQIIdrwRR2CTGg5CWq9c-mdeM8KlW52B7w2NkiQ9Lpvh8lcbx7gETJy_Jo1JEy3ImnJUXmrTRz-jA1dLEYmSscxqLAQjhdnfbTOuS3_JOEhaRLEzm-j'
    }
  ];

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'done':
        return (
          <span className="text-base fill-current" style={{ fontVariationSettings: "'FILL' 1" }}>
            ✓
          </span>
        );
      case 'review':
        return (
          <span className="text-base fill-current" style={{ fontVariationSettings: "'FILL' 1" }}>
            ⚠
          </span>
        );
      case 'extracting':
        return (
          <span className="text-base animate-spin">
            ↻
          </span>
        );
      default:
        return null;
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'done':
        return 'text-status-green';
      case 'review':
        return 'text-status-yellow';
      case 'extracting':
        return 'text-status-blue';
      default:
        return 'text-gray-500';
    }
  };

  return (
    <div className="relative mx-auto flex h-screen max-h-[960px] w-full max-w-[480px] flex-col overflow-hidden bg-background-light dark:bg-background-dark font-display text-text-light dark:text-text-dark">
      {/* Header */}
      <header className="flex shrink-0 items-center justify-between border-b border-border-light dark:border-border-dark bg-surface-light dark:bg-surface-dark px-4 py-3">
        <button className="flex size-10 cursor-pointer items-center justify-center rounded-full">
          <span className="material-symbols-outlined text-2xl">←</span>
        </button>
        <h1 className="text-lg font-semibold">Scanned Invoices</h1>
        <button className="flex size-10 cursor-pointer items-center justify-center rounded-full">
          <span className="material-symbols-outlined text-2xl">🔍</span>
        </button>
      </header>

      {/* Main Content */}
      <main className="flex-grow overflow-y-auto p-4">
        <div className="flex flex-col gap-4">
          {invoices.map((invoice) => (
            <button
              key={invoice.id}
              className="flex items-center gap-4 rounded-xl bg-surface-light p-3 shadow-sm transition-shadow hover:shadow-md dark:bg-surface-dark text-left w-full"
              onClick={() => navigate('/extract-invoice-data')}
            >
              <img 
                alt="Invoice thumbnail" 
                className="h-16 w-16 flex-shrink-0 rounded-lg object-cover bg-gray-200" 
                src={invoice.thumbnail}
              />
              <div className="flex-grow">
                <h2 className="font-semibold text-text-light dark:text-text-dark">{invoice.title}</h2>
                <div className={`mt-1 flex items-center gap-1.5 text-sm ${getStatusColor(invoice.status)}`}>
                  {getStatusIcon(invoice.status)}
                  <p>{invoice.statusText}</p>
                </div>
              </div>
              <span className="material-symbols-outlined text-text-light-secondary dark:text-text-dark-secondary">›</span>
            </button>
          ))}
        </div>
      </main>

      {/* Floating Action Button */}
      <div className="absolute bottom-6 right-6">
        <button 
          className="flex h-16 w-16 items-center justify-center rounded-full bg-primary text-white shadow-lg transition-transform hover:scale-105 active:scale-95"
          onClick={() => navigate('/take-picture')}
        >
          <span className="text-4xl" style={{ fontVariationSettings: "'FILL' 1" }}>📷</span>
        </button>
      </div>
    </div>
  );
}