import { Component, type ReactNode } from 'react';

type Props = {
  children: ReactNode;
};

type State = {
  hasError: boolean;
  message?: string;
};

export default class ErrorBoundary extends Component<Props, State> {
  state: State = { hasError: false };

  static getDerivedStateFromError(error: unknown): State {
    const message = error instanceof Error ? error.message : 'Unknown error';
    return { hasError: true, message };
  }

  componentDidCatch(error: unknown) {
    // eslint-disable-next-line no-console
    console.error('UI crashed:', error);
  }

  render() {
    if (this.state.hasError) {
      return (
        <div className="rounded-3xl bg-white border border-border p-8 text-center shadow-sm space-y-4">
          <div className="text-lg font-bold text-text-main">Bir hata oluştu</div>
          <div className="text-sm text-text-muted">{this.state.message || 'Beklenmeyen bir hata oluştu.'}</div>
          <button
            type="button"
            onClick={() => window.location.assign('/')}
            className="inline-flex items-center justify-center rounded-2xl bg-primary hover:bg-primary-hover text-white px-5 py-2.5 text-sm font-semibold transition-colors"
          >
            Ana sayfaya dön
          </button>
        </div>
      );
    }

    return this.props.children;
  }
}
