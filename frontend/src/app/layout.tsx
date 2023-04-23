import '../css/global.css'
import { Inter } from 'next/font/google'
import { getServerSession } from "next-auth"
import { authOptions } from "@/pages/api/auth/[...nextauth]"
import { SessionProvider } from "@/components/SessionProvider";
import Login from '@/components/Login';

const inter = Inter({ subsets: ['latin'] })

export const metadata = {
  title: 'Play With ChatGPT',
  description: 'Playground for ChatGPT development',
}

export default async function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  const session = await getServerSession(authOptions);

  return (
    <html lang="en">
      <body className={inter.className}>

      <SessionProvider session={session}>
        {
          !session ? (
            <Login />
          ) : (
              <div className="flex">
                <div className="bg-[#343541] flex-1">{children}</div>
              </div>

          )}
      </SessionProvider>
      </body>

    </html>
  )
}
