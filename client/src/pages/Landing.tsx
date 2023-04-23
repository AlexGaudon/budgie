import budgieImage from "/bird.svg"; // Replace with your actual image path

const LandingPage = () => {
    return (
        <div className="text-white min-h-screen">
            <div className="container mx-auto px-4 py-12">
                <div className="flex flex-wrap">
                    <div className="w-full md:w-1/2">
                        <img
                            src={budgieImage}
                            alt="Budgie App"
                            className="mx-auto max-w-sm"
                        />
                    </div>
                    <div className="w-full md:w-1/2 mt-8 md:mt-0 md:pl-8">
                        <h1 className="text-4xl font-semibold mb-6">
                            Welcome to Budgie!
                        </h1>
                        <h2 className="text-2xl font-medium mb-8">
                            Budgeting and Personal Finance Made Easy
                        </h2>
                        <p className="mb-4">
                            Budgie is a powerful and user-friendly app that
                            helps you take control of your finances. With
                            Budgie, you can effortlessly log transactions,
                            categorize them, and set budgets to manage your
                            expenses with ease.
                        </p>
                        <ul className="list-disc list-inside mb-8">
                            <li>Effortless Transaction Logging</li>
                            <li>Smart Categorization</li>
                            <li>Budget Tracking</li>
                            <li>Insightful Reports</li>
                            <li>Secure and Private</li>
                        </ul>
                        <div className="flex justify-center md:justify-start">
                            <button
                                className="bg-blue-500 hover:bg-blue-600 text-white font-semibold py-3 px-6 rounded"
                                onClick={() => {
                                    // Handle sign up action
                                }}
                            >
                                Sign Up Now
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default LandingPage;
