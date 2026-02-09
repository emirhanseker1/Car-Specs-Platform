import { useState, useRef, useEffect } from 'react';
import { Send, Sparkles, Zap, RefreshCw, X } from 'lucide-react';
import { api } from '../../services/api';

interface Message {
    id: string;
    text: string;
    sender: 'user' | 'bot';
    timestamp: Date;
}



export default function ChatWidget() {
    const [isOpen, setIsOpen] = useState(false);
    const [messages, setMessages] = useState<Message[]>([
        {
            id: '1',
            text: 'Merhaba! Ben Car Spec AI Asistan. Araçlarla ilgili aklınıza takılan her şeyi sorabilirsiniz.',
            sender: 'bot',
            timestamp: new Date(),
        },
    ]);
    const [inputText, setInputText] = useState('');
    const [isLoading, setIsLoading] = useState(false);
    const messagesEndRef = useRef<HTMLDivElement>(null);
    const inputRef = useRef<HTMLInputElement>(null);

    const toggleChat = () => {
        setIsOpen(!isOpen);
        if (!isOpen) {
            setTimeout(() => inputRef.current?.focus(), 100);
        }
    };

    const handleNewChat = () => {
        setMessages([
            {
                id: Date.now().toString(),
                text: 'Merhaba! Yeni bir sayfa açtık. Araçlarla ilgili aklınıza takılan her şeyi sorabilirsiniz.',
                sender: 'bot',
                timestamp: new Date(),
            },
        ]);
        setInputText('');
        inputRef.current?.focus();
    };

    const scrollToBottom = () => {
        messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
    };

    useEffect(() => {
        if (isOpen) scrollToBottom();
    }, [messages, isOpen]);

    useEffect(() => {
        const handleKeyDown = (e: KeyboardEvent) => {
            if (e.key === 'Escape' && isOpen) {
                setIsOpen(false);
            }
            if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
                e.preventDefault();
                toggleChat();
            }
        };

        const handleCustomOpen = () => {
            setIsOpen(true);
            setTimeout(() => inputRef.current?.focus(), 100);
        };

        window.addEventListener('keydown', handleKeyDown);
        window.addEventListener('open-chat-trigger', handleCustomOpen);

        return () => {
            window.removeEventListener('keydown', handleKeyDown);
            window.removeEventListener('open-chat-trigger', handleCustomOpen);
        };
    }, [isOpen]);

    const handleSendMessage = async (textOverride?: string) => {
        const textToSend = textOverride || inputText;
        if (!textToSend.trim() || isLoading) return;

        const userMessage: Message = {
            id: Date.now().toString(),
            text: textToSend,
            sender: 'user',
            timestamp: new Date(),
        };

        setMessages((prev) => [...prev, userMessage]);
        setInputText('');
        setIsLoading(true);

        try {
            const response = await api.sendChatMessage(userMessage.text);

            const botMessage: Message = {
                id: (Date.now() + 1).toString(),
                text: response.response || 'Bir hata oluştu.',
                sender: 'bot',
                timestamp: new Date(),
            };

            setMessages((prev) => [...prev, botMessage]);
        } catch (error) {
            const errorMessage: Message = {
                id: (Date.now() + 1).toString(),
                text: error instanceof Error ? error.message : 'Üzgünüm, bağlantı hatası oluştu. Lütfen tekrar deneyin.',
                sender: 'bot',
                timestamp: new Date(),
            };
            setMessages((prev) => [...prev, errorMessage]);
        } finally {
            setIsLoading(false);
        }
    };

    const handleKeyPress = (e: React.KeyboardEvent) => {
        if (e.key === 'Enter') {
            handleSendMessage();
        }
    };

    return (
        <>
            {/* Floating Trigger Button */}
            {!isOpen && (
                <button
                    onClick={toggleChat}
                    className="fixed bottom-8 right-8 z-50 group flex items-center justify-center transition-all duration-300 hover:scale-105"
                >
                    <div className="absolute inset-0 bg-primary/30 rounded-full animate-ping opacity-75"></div>
                    <div className="relative bg-zinc-900 border border-primary/50 text-white p-4 rounded-full shadow-[0_0_30px_rgba(255,90,31,0.3)]">
                        <Sparkles className="w-6 h-6 text-primary animate-pulse" />
                    </div>
                </button>
            )}

            {/* Modal Overlay */}
            {isOpen && (
                <div className="fixed inset-0 z-50 flex items-center justify-center p-4 sm:p-6 font-sans">
                    {/* Backdrop */}
                    <div
                        className="absolute inset-0 bg-black/80 backdrop-blur-sm transition-opacity duration-300 animate-in fade-in"
                        onClick={() => setIsOpen(false)}
                    />

                    {/* Main Modal Container - Sahibinden Style Structure */}
                    <div className="relative w-full max-w-4xl h-[85vh] bg-zinc-950 border border-zinc-800 rounded-xl shadow-2xl flex flex-col overflow-hidden animate-in zoom-in-95 duration-200">

                        {/* 1. Header */}
                        <div className="flex items-center justify-between px-6 py-4 border-b border-zinc-800 bg-zinc-900/50">
                            <div className="flex items-center gap-3">
                                <div className="w-8 h-8 rounded-lg bg-primary/10 flex items-center justify-center border border-primary/20">
                                    <Sparkles className="w-5 h-5 text-primary" />
                                </div>
                                <div>
                                    <h3 className="text-lg font-bold text-white leading-none">Car Spec AI Asistan</h3>
                                    <span className="text-xs text-zinc-400">Powered by LLM</span>
                                </div>
                            </div>

                            <div className="flex items-center gap-2">
                                <button
                                    onClick={handleNewChat}
                                    className="flex items-center gap-2 px-3 py-1.5 rounded-lg bg-zinc-800 hover:bg-zinc-700 text-zinc-300 text-sm transition-colors border border-white/5"
                                    title="Sohbeti Temizle"
                                >
                                    <RefreshCw className="w-4 h-4" />
                                    <span className="hidden sm:inline">Yeni Sohbet</span>
                                </button>
                                <button
                                    onClick={() => setIsOpen(false)}
                                    className="p-2 rounded-lg hover:bg-red-500/10 hover:text-red-500 text-zinc-400 transition-colors"
                                >
                                    <X className="w-5 h-5" />
                                </button>
                            </div>
                        </div>

                        {/* 2. Messages Area */}
                        <div className="flex-1 overflow-y-auto p-6 space-y-6 bg-zinc-950/50 scrollbar-thin scrollbar-thumb-zinc-800 scrollbar-track-transparent">
                            {messages.map((msg) => (
                                <div key={msg.id} className={`flex gap-4 ${msg.sender === 'user' ? 'justify-end' : 'justify-start'}`}>
                                    {msg.sender === 'bot' && (
                                        <div className="w-8 h-8 rounded-full bg-zinc-800 border border-white/5 flex items-center justify-center flex-shrink-0 mt-1">
                                            <Sparkles className="w-4 h-4 text-primary" />
                                        </div>
                                    )}

                                    <div className={`max-w-[85%] space-y-1 ${msg.sender === 'user' ? 'items-end flex flex-col' : ''}`}>
                                        <div className={`px-5 py-3.5 rounded-2xl text-[15px] leading-relaxed shadow-sm whitespace-pre-wrap ${msg.sender === 'user'
                                            ? 'bg-primary text-white rounded-br-sm'
                                            : 'bg-zinc-900 border border-zinc-800 text-gray-200 rounded-bl-sm'
                                            }`}>
                                            {msg.text}
                                        </div>
                                    </div>
                                </div>
                            ))}

                            {isLoading && (
                                <div className="flex gap-4 justify-start">
                                    <div className="w-8 h-8 rounded-full bg-zinc-800 flex items-center justify-center flex-shrink-0">
                                        <Zap className="w-4 h-4 text-primary animate-pulse" />
                                    </div>
                                    <div className="flex items-center gap-1.5 h-10 px-4 bg-zinc-900 rounded-2xl border border-zinc-800">
                                        <div className="w-2 h-2 bg-primary/50 rounded-full animate-bounce" style={{ animationDelay: '0ms' }} />
                                        <div className="w-2 h-2 bg-primary/50 rounded-full animate-bounce" style={{ animationDelay: '150ms' }} />
                                        <div className="w-2 h-2 bg-primary/50 rounded-full animate-bounce" style={{ animationDelay: '300ms' }} />
                                    </div>
                                </div>
                            )}
                            <div ref={messagesEndRef} />
                        </div>

                        {/* 3. Footer Area */}
                        <div className="bg-zinc-900 border-t border-zinc-800 p-4">



                            {/* Input Field */}
                            <div className="relative">
                                <input
                                    ref={inputRef}
                                    type="text"
                                    value={inputText}
                                    onChange={(e) => setInputText(e.target.value)}
                                    onKeyDown={handleKeyPress}
                                    placeholder="Aradığın araç hakkında detay sor... (Örn: Golf 7 1.6 TDI yakıt tüketimi)"
                                    className="w-full bg-zinc-950 border border-zinc-700 focus:border-primary/50 text-white placeholder-zinc-500 rounded-xl py-4 pl-5 pr-14 text-base outline-none transition-all focus:ring-1 focus:ring-primary/30"
                                    autoFocus
                                    disabled={isLoading}
                                />
                                <button
                                    onClick={() => handleSendMessage()}
                                    disabled={!inputText.trim() || isLoading}
                                    className="absolute right-2 top-1/2 -translate-y-1/2 p-2 bg-primary hover:bg-primary-hover text-white rounded-lg disabled:opacity-50 disabled:bg-zinc-700 transition-all"
                                >
                                    <Send className="w-5 h-5" />
                                </button>
                            </div>

                            {/* Disclaimer */}
                            <div className="text-center mt-3">
                                <p className="text-xs text-zinc-500 flex items-center justify-center gap-1.5">
                                    <Sparkles className="w-3 h-3" />
                                    Car Spec AI hata yapabilir. Önemli bilgileri kontrol ediniz.
                                </p>
                            </div>
                        </div>

                    </div>
                </div>
            )}
        </>
    );
}
