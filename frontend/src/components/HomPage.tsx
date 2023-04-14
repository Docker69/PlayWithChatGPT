import React from 'react';
import Layout from './Layout';
import ChatSession from './ChatSession';

const HomePage: React.FC = () => {
  return (
    <Layout>
       <ChatSession />
    </Layout>
  );
};

export default HomePage;