import { SignUp as ClerkSignUp, useUser } from '@clerk/clerk-react';
import { Navigate } from 'react-router-dom';
import { PocketLogo } from '../components/PocketLogo';

export const SignUp = () => {
  const { isSignedIn, isLoaded } = useUser();

  if (!isLoaded) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="text-gray-600">Loading...</div>
      </div>
    );
  }

  if (isSignedIn) {
    return <Navigate to="/" replace />;
  }

  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-gray-50 px-4">
      <div className="mb-8 flex flex-col items-center gap-4">
        <PocketLogo className="text-blue-600" size={64} />
        <h1 className="text-3xl font-bold text-gray-900">Mi Bolsillo</h1>
      </div>
      <ClerkSignUp
        routing="virtual"
        signInUrl="/login"
        afterSignUpUrl="/"
        redirectUrl="/"
      />
    </div>
  );
};
